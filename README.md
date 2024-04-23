# my-umich-events

## How to run 
In development:

```
uvicorn main:app --reload
```

## Ideas

### Overall design
The website will only have 3 paths:
- `/`: homepage that describes the project, etc and displays most popular events of the week
- `/recs?user_id=`: shows recommended events for user
- `/settings?user_id=`: user inputs their preferences

### Home view
- Site description
- Popular events (by adds to gCalendar)

### Recommendations view
List of recommended events
- why it was recommended (keyword, ai, others seem to enjoy, etc.)
- expandable description
- button to add to gCalendar (and log positive feedback)
- button to show less

### Preferences view
- Select program
- Input schedule (i.e. available time slots)
- research interests (str - maybe further along)
- keywords
- researcher talks w/h index above threshold

## Where to apply ML?
- Summarize card titles and descriptions
- Recommendations 