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
  "timeslots": [37, 38, 39, 40, 67, 68, 69, 70, 113, 114, 115, 116]
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
  "timeslots": [35, 36, 37, 38, 39, 40, 41, 42, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 173, 174, 175, 176, 177, 178, 179, 180]
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
  "timeslots": [39, 40, 41, 42, 67, 68, 69, 70, 91, 92, 93, 94]
}
```

## Why These Are Compatible:

### 1. Same Activity ID
- All three requests have `activity_id: 5` (tennis)

### 2. Shared Week Timeslots
Let's calculate the week_timeslots for each (half-hour intervals):

**Request 1 week_timeslots:** [37, 38, 39, 40, 67, 68, 69, 70, 113, 114, 115, 116] 
- Monday 18-20: [37, 38, 39, 40] (18*2+1 to 20*2)
- Wednesday 18-20: [67, 68, 69, 70] (48+18*2+1 to 48+20*2)  
- Friday 17-19: [113, 114, 115, 116] (96+17*2+1 to 96+19*2)

**Request 2 week_timeslots:** [35, 36, 37, 38, 39, 40, 41, 42, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 173, 174, 175, 176, 177, 178, 179, 180]
- Monday 17-21: [35, 36, 37, 38, 39, 40, 41, 42] (17*2+1 to 21*2)
- Wednesday 16-22: [65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92] (48+16*2+1 to 48+22*2)
- Saturday 14-18: [173, 174, 175, 176, 177, 178, 179, 180] (144+14*2+1 to 144+18*2)

**Request 3 week_timeslots:** [39, 40, 41, 42, 67, 68, 69, 70, 91, 92, 93, 94]
- Monday 19-21: [39, 40, 41, 42] (19*2+1 to 21*2)
- Wednesday 18-20: [67, 68, 69, 70] (48+18*2+1 to 48+20*2)
- Thursday 18-20: [91, 92, 93, 94] (72+18*2+1 to 72+20*2)

**Shared timeslots between all three:** [39, 40, 67, 68, 69, 70]
- Monday 19:00-20:00 (timeslots 39, 40)
- Wednesday 18:00-20:00 (timeslots 67, 68, 69, 70)

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