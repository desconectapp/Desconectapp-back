# Example: Three Compatible Activity Requests

Here are three activity requests that are all compatible and should create a group:

## Activity Request 1 (User ID: 101)
```json
{
  "user_id": 1,
  "activity_id": 5,
  "description": "Looking for tennis partners in Palermo",
  "longitude": -58.4173,
  "latitude": -34.5755,
  "search_radius": 5,
  "max_participants": 6,
  "participants_needed": 3,
  "schedules": {
    "monday": [{"start": 18, "end": 20}],
    "wednesday": [{"start": 18, "end": 20}],
    "friday": [{"start": 17, "end": 19}]
  }
}
```

## Activity Request 2 (User ID: 102)
```json
{
  "user_id": 2,
  "activity_id": 5,
  "description": "Tennis group in Palermo area",
  "longitude": -58.4200,
  "latitude": -34.5780,
  "search_radius": 4,
  "max_participants": 8,
  "participants_needed": 2,
  "schedules": {
    "monday": [{"start": 17, "end": 21}],
    "wednesday": [{"start": 16, "end": 22}],
    "saturday": [{"start": 14, "end": 18}]
  }
}
```

## Activity Request 3 (User ID: 103)
```json
{
  "user_id": 3,
  "activity_id": 5,
  "description": "Tennis matches near Palermo",
  "longitude": -58.4150,
  "latitude": -34.5730,
  "search_radius": 6,
  "max_participants": 4,
  "participants_needed": 3,
  "schedules": {
    "monday": [{"start": 19, "end": 21}],
    "wednesday": [{"start": 18, "end": 20}],
    "thursday": [{"start": 18, "end": 20}]
  }
}
```

## Why These Are Compatible:

### 1. Same Activity ID
- All three requests have `activity_id: 5` (tennis)

### 2. Shared Week Hours
Let's calculate the week_hours for each:

**Request 1 week_hours:** [18, 19, 66, 67, 113, 114] 
- Monday 18-20: [18, 19]
- Wednesday 18-20: [48+18, 48+19] = [66, 67]  
- Friday 17-19: [96+17, 96+18] = [113, 114]

**Request 2 week_hours:** [17, 18, 19, 20, 64, 65, 66, 67, 68, 69, 134, 135, 136, 137]
- Monday 17-21: [17, 18, 19, 20]
- Wednesday 16-22: [64, 65, 66, 67, 68, 69]
- Saturday 14-18: [134, 135, 136, 137]

**Request 3 week_hours:** [19, 20, 66, 67, 90, 91]
- Monday 19-21: [19, 20]
- Wednesday 18-20: [66, 67]
- Thursday 18-20: [90, 91]

**Shared hours between all three:** [19, 66, 67]
- Monday 19:00 (hour 19)
- Wednesday 18:00-19:00 (hours 66, 67)

### 3. Location Compatibility
All locations are in Palermo, Buenos Aires:
- Request 1: (-58.4173, -34.5755) with 5km radius
- Request 2: (-58.4200, -34.5780) with 4km radius  
- Request 3: (-58.4150, -34.5730) with 6km radius

Distance between any two points is ~0.5-0.8km, well within all search radii.

### 4. Participant Count Compatibility
- Request 1: needs 3, max 6
- Request 2: needs 2, max 8  
- Request 3: needs 2, max 4

Each request's max_participants (6, 8, 4) is >= others' participants_needed (3, 2, 2) ✓

## Expected Group Formation:

When these three requests are processed:

1. **First match** (Req 1 + Req 2) creates partial match with:
   - Week hours: [18, 19, 66, 67] (intersection)
   - Location: midpoint (-58.41865, -34.57675)
   - Search radius: 4.5km (average)
   - Min participants: 3 (highest of 3, 2)
   - Max participants: 6 (lowest of 6, 8)

2. **Second match** (Partial + Req 3) has enough members (3 total) to create a group:
   - Final week hours: [19, 66, 67] (intersection of all three)
   - Final location: midpoint of all three locations
   - Members: Users 101, 102, 103 