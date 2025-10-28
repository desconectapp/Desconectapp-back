-- name: CreateCommunity :one
WITH inserted_community AS (
  INSERT INTO communities (name, avatar_url, location, location_name, activity_id, week_timeslots, description)
  VALUES ($1, $2, $3, $4, $5, $6, $7)
  RETURNING *
), inserted_members AS (
  INSERT INTO communities_members (user_id, community_id, is_admin)
  SELECT 
    u.id, 
    c.id, 
    u.id = ANY(sqlc.arg(admin_user_ids)::int[]) AS is_admin
  FROM users u
  CROSS JOIN inserted_community c
  WHERE u.id = ANY(sqlc.arg(user_ids)::int[])
  ON CONFLICT DO NOTHING
  RETURNING user_id, community_id, is_admin
)
SELECT 
  c.id,
  c.name,
  c.avatar_url,
  c.location,
  c.location_name,
  c.activity_id,
  c.week_timeslots,
  c.description,
  c.created_at,
  CAST(
    COALESCE(
      array_agg(m.user_id) FILTER (WHERE m.user_id IS NOT NULL),
      ARRAY[]::int[]
    ) AS int[]
  ) AS members,
  CAST(
    COALESCE(
      array_agg(m.user_id) FILTER (WHERE m.is_admin = true),
      ARRAY[]::int[]
    ) AS int[]
  ) AS admin_members,
  a.name AS activity_name,
  a.icon AS activity_icon
FROM inserted_community c
LEFT JOIN inserted_members m ON c.id = m.community_id
JOIN activities a ON c.activity_id = a.id
GROUP BY c.id, c.name, c.avatar_url, c.description, c.location, c.location_name, 
         c.activity_id, c.week_timeslots, c.created_at, 
         a.name, a.icon;

-- name: ListCommunities :many
WITH selected_communities AS (
  SELECT *
  FROM communities
  ORDER BY created_at DESC
  LIMIT $1 OFFSET $2
)
SELECT 
  c.id,
  c.name,
  c.avatar_url,
  c.location_name,
  c.description,
  c.created_at,
  a.name AS activity,
  a.icon,
  COUNT(cm.user_id) AS members_count
FROM selected_communities c
JOIN activities a ON c.activity_id = a.id
LEFT JOIN communities_members cm ON c.id = cm.community_id
GROUP BY c.id, c.name, c.avatar_url, c.location_name, c.description, c.created_at, a.name, a.icon
ORDER BY c.created_at DESC;

-- name: GetCommunity :one
SELECT 
  c.id,
  c.name,
  c.avatar_url,
  c.location_name,
  c.location,
  c.week_timeslots,
  c.description,
  c.created_at,
  a.name AS activity,
  a.icon,
  CAST(
    COALESCE(
      json_agg(
        DISTINCT jsonb_build_object(
          'user_id', u.id,
          'uuid', u.uuid,
          'is_admin', cm.is_admin
        )
      ) FILTER (WHERE u.id IS NOT NULL),
      '[]'
    ) AS json
  ) AS members
FROM communities c
JOIN activities a ON c.activity_id = a.id
LEFT JOIN communities_members cm ON cm.community_id = c.id
LEFT JOIN users u ON cm.user_id = u.id
WHERE c.id = $1
GROUP BY c.id, c.name, c.avatar_url, c.location_name, c.location, c.description, c.week_timeslots, c.created_at, a.name, a.icon;

-- name: DeleteCommunity :one
DELETE FROM communities
WHERE id = $1
RETURNING id;

-- name: AddUserToCommunity :exec
INSERT INTO communities_members (community_id, user_id, is_admin)
VALUES ($1, $2, COALESCE($3, false))
ON CONFLICT DO NOTHING;

-- name: BatchAddUsersToCommunity :exec
INSERT INTO communities_members (user_id, community_id, is_admin)
SELECT id, sqlc.arg(community_id), false
FROM users
WHERE id = ANY(sqlc.arg(user_ids)::int[])
ON CONFLICT DO NOTHING;

-- name: GetCommunityMembers :many
SELECT 
  u.id, 
  u.uuid, 
  p.name, 
  p.avatar_url,
  cm.is_admin
FROM users u
JOIN profiles p ON u.id = p.user_id
JOIN communities_members cm ON u.id = cm.user_id
WHERE cm.community_id = $1;

-- name: ListUserCommunities :many
WITH user_communities AS (
  SELECT c.*
  FROM communities c
  JOIN communities_members cm ON c.id = cm.community_id
  WHERE cm.user_id = $3
  ORDER BY c.created_at DESC
  LIMIT $1 OFFSET $2
)
SELECT 
  c.id,
  c.name,
  c.avatar_url,
  c.location_name,
  c.location,
  c.week_timeslots,
  c.description,
  c.created_at,
  a.name AS activity,
  a.icon,
  COUNT(DISTINCT cm_all.user_id) AS members_count,
  COALESCE(
    json_agg(
      DISTINCT jsonb_build_object(
        'uuid', u.uuid,
        'name', p.name,
        'is_admin', cm_all.is_admin
      )
    ) FILTER (WHERE u.id IS NOT NULL),
    '[]'
  ) AS members
FROM user_communities c
JOIN activities a ON c.activity_id = a.id
LEFT JOIN communities_members cm_all ON c.id = cm_all.community_id
LEFT JOIN users u ON cm_all.user_id = u.id
LEFT JOIN profiles p ON cm_all.user_id = p.user_id
GROUP BY c.id, c.name, c.avatar_url, c.description, c.location_name, c.location, c.week_timeslots, c.created_at, a.name, a.icon
ORDER BY c.created_at DESC;

-- name: ExitCommunity :exec
DELETE FROM communities_members
WHERE community_id = $1 AND user_id = $2;

-- name: UpdateCommunityDescription :exec
UPDATE communities
SET description = $2
WHERE id = $1;

-- name: ChangeCommunityName :exec
UPDATE communities
SET name = $2
WHERE id = $1;

-- name: ChangeCommunityLocation :exec
UPDATE communities
SET location = $2, location_name = $3
WHERE id = $1;

-- name: UpdateCommunityAvatar :exec
UPDATE communities
SET avatar_url = $2
WHERE id = $1;

-- name: PromoteUserToAdmin :exec
UPDATE communities_members
SET is_admin = true
WHERE community_id = $1 AND user_id = $2;

-- name: DemoteAdminToUser :exec
UPDATE communities_members
SET is_admin = false
WHERE community_id = $1 AND user_id = $2;

-- name: GetCommunitiesWithLocation :many
SELECT 
  c.id,
  c.name,
  c.avatar_url,
  c.description,
  c.location,
  c.location_name,
  c.week_timeslots,
  a.name AS activity_name,
  a.icon,
  COUNT(cm.user_id) AS member_count,
  CAST((6371 * acos(
      cos(radians(sqlc.arg('latitude')::float)) * 
      cos(radians(CAST(split_part(c.location, ',', 1) AS float))) *
      cos(radians(CAST(split_part(c.location, ',', 2) AS float)) - radians(sqlc.arg('longitude')::float)) +
      sin(radians(sqlc.arg('latitude')::float)) * 
      sin(radians(CAST(split_part(c.location, ',', 1) AS float)))
  )) AS float) AS distance_km
FROM communities c
JOIN activities a ON c.activity_id = a.id
LEFT JOIN communities_members cm ON c.id = cm.community_id
WHERE c.location IS NOT NULL
  AND c.location != ''
  AND (sqlc.narg('activity_id')::int IS NULL OR c.activity_id = sqlc.narg('activity_id')::int)
  AND (6371 * acos(
        cos(radians(sqlc.arg('latitude')::float)) * 
        cos(radians(CAST(split_part(c.location, ',', 1) AS float))) *
        cos(radians(CAST(split_part(c.location, ',', 2) AS float)) - radians(sqlc.arg('longitude')::float)) +
        sin(radians(sqlc.arg('latitude')::float)) * 
        sin(radians(CAST(split_part(c.location, ',', 1) AS float)))
    )) <= sqlc.arg('radius')::float
GROUP BY c.id, c.name, c.avatar_url, c.location, c.description, c.location_name, c.week_timeslots, a.name, a.icon
ORDER BY distance_km
LIMIT $1 OFFSET $2;

-- name: GetPreferredCommunities :many
SELECT 
  c.id,
  c.name,
  c.avatar_url,
  c.location,
  c.location_name,
  c.week_timeslots,
  c.description,
  c.activity_id,
  c.created_at,
  COUNT(cm.user_id) AS member_count,
  a.name AS activity_name,
  a.icon AS activity_icon
FROM communities c
JOIN users_preference up ON c.activity_id = up.activity_id
LEFT JOIN communities_members cm ON c.id = cm.community_id
JOIN activities a ON c.activity_id = a.id
WHERE up.user_id = $1
  AND NOT EXISTS (
      SELECT 1
      FROM communities_members cm2
      WHERE cm2.community_id = c.id
        AND cm2.user_id = $1
  )
GROUP BY c.id, c.name, c.avatar_url, c.location, c.location_name, c.description, c.week_timeslots, c.activity_id, c.created_at, a.name, a.icon
ORDER BY member_count ASC, c.created_at DESC
LIMIT $2 OFFSET $3;

