from typing import Optional, Union

from fastapi import FastAPI, Form
from fastapi.staticfiles import StaticFiles
from fastapi.responses import HTMLResponse
from pydantic import BaseModel

import json, time

app = FastAPI()
app.mount("/static", StaticFiles(directory = "static"), name = "static")


@app.get("/")
async def root():
    
    with open('data/events.json', 'r') as f:
        events = json.load(f)

    with open('templates/card.html', 'r') as f:
        card_template = f.read()

    with open("templates/index.html") as f:
        index_template = f.read()

    # cut off text to l chars and add '...'
    cutoff = lambda s, l: s if len(s) <= l else s[:l] + '...'

    # cutoff time
    cutoff_time = lambda s: ''.join([i for i in s.split(':')[:1] if i != '00'])

    time_str = lambda s, e: f'{cutoff_time(s)} - {cutoff_time(e)}'
    
    cards_html = []
    for event in events:
        try:
            card = card_template[::]
            card = card.replace('{{title}}', cutoff(event['event_title'], 50))
            time_s = time_str(event['time_start'], event['time_end'])
            card = card.replace('{{subtitle}}', f'{event["event_type"].split("/")[0]} | {time_s} | {cutoff(event["location_name"], 20)}')
            card = card.replace('{{description}}', cutoff(event['description'], 200))
            card = card.replace('{{event_link}}', event['permalink'])
            card = card.replace('{{event_link}}', event['permalink'])
            cards_html.append(card)
        except:
            print('Could not parse event', event['event_title'])
    
    
    index_template = index_template.replace('{{events}}', '\n\t\t\t'.join(cards_html))

    return HTMLResponse(content = index_template, status_code = 200)


@app.get("/prefs", response_class = HTMLResponse)
async def get_prefs():
    # Read the template file
    with open("templates/prefs.html") as f:
        template = f.read()

    with open('static/json/programs.json', 'r') as f:
        programs = json.load(f)
    
    programs_HTML = []
    for program in programs:
        programs_HTML.append(f'<option value="{program}">{program}</option>')

    template = template.replace('{{programs}}', '\n\t\t\t'.join(programs_HTML))

    eventTyes = ['Seminars', 'Sports', 'Social']
    eventTypes_HTML = []
    for eventType in eventTyes:
        eventTypes_HTML.append(
            f'''<div class="form-check">
                <input name="inc{eventType}" value="true" hx-put="/prefs" hx-indicator="#eventTypes_indicator" class="form-check-input" type="checkbox" value="" id="flexCheckDefault">
                <label class="form-check-label" for="flexCheckDefault">
                    {eventType}
                </label>
            </div>'''
        )
    template = template.replace('{{eventCheckboxes}}', '\n\t\t\t'.join(eventTypes_HTML))
    
    return HTMLResponse(content = template, status_code = 200)

@app.put("/prefs")
async def set_prefs(
        educationLevel: Optional[str] = Form(None),
        program: Optional[str] = Form(None),
        campusLocation: Optional[str] = Form(None),
        researchInterests: Optional[str] = Form(None),
        incSeminars: Optional[Union[bool, None]] = Form(None),
        incSports: Optional[Union[bool, None]] = Form(None),
        incSocial: Optional[Union[bool, None]] = Form(None),
        keywordsToAvoid: Optional[str] = Form(None),
    ):

    # Can only change a single param at a time
    changed_param = {k:v for k, v in locals().items() if v is not None}
    if len(changed_param) != 1:
        print('Invalid params')
        return HTMLResponse(content = 'Invalid params', status_code = 400)

    print(changed_param)
    time.sleep(1)
    return HTMLResponse(content = 'Saved', status_code = 200)

@app.post("/upvote")
async def upvote(event_id: str = None):
    time.sleep(1)
    return HTMLResponse(content = 'Upvoted', status_code = 200)