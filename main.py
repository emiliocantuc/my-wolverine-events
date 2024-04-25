from typing import Optional, Union

from fastapi import FastAPI, Form
from fastapi.staticfiles import StaticFiles
from fastapi.responses import HTMLResponse
from fastapi.templating import Jinja2Templates

from pydantic import BaseModel

import json, time

app = FastAPI()
app.mount("/static", StaticFiles(directory = "static"), name = "static")
templates = Jinja2Templates(directory="templates")


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
            for _ in range(3):
                card = card.replace('{{event_id}}', event['id'])
            card = card.replace('{{title}}', cutoff(event['event_title'], 50))
            card = card.replace('{{calendar_link}}', event['gcal_link'])
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
        interests: Optional[str] = Form(None),
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


@app.put("/addcal/{event_id}")
async def addcal(event_id: str = None):
    time.sleep(1)
    filled_cal = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-calendar-check-fill" viewBox="0 0 16 16"><path d="M4 .5a.5.5 0 0 0-1 0V1H2a2 2 0 0 0-2 2v1h16V3a2 2 0 0 0-2-2h-1V.5a.5.5 0 0 0-1 0V1H4zM16 14V5H0v9a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2m-5.146-5.146-3 3a.5.5 0 0 1-.708 0l-1.5-1.5a.5.5 0 0 1 .708-.708L7.5 10.793l2.646-2.647a.5.5 0 0 1 .708.708"/></svg>'
    return HTMLResponse(content = filled_cal, status_code = 200)

@app.put("/upvote/{event_id}")
async def upvote(event_id: str = None):
    time.sleep(1)
    thumb_up_filled = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-hand-thumbs-up-fill" viewBox="0 0 16 16"><path d="M6.956 1.745C7.021.81 7.908.087 8.864.325l.261.066c.463.116.874.456 1.012.965.22.816.533 2.511.062 4.51a10 10 0 0 1 .443-.051c.713-.065 1.669-.072 2.516.21.518.173.994.681 1.2 1.273.184.532.16 1.162-.234 1.733q.086.18.138.363c.077.27.113.567.113.856s-.036.586-.113.856c-.039.135-.09.273-.16.404.169.387.107.819-.003 1.148a3.2 3.2 0 0 1-.488.901c.054.152.076.312.076.465 0 .305-.089.625-.253.912C13.1 15.522 12.437 16 11.5 16H8c-.605 0-1.07-.081-1.466-.218a4.8 4.8 0 0 1-.97-.484l-.048-.03c-.504-.307-.999-.609-2.068-.722C2.682 14.464 2 13.846 2 13V9c0-.85.685-1.432 1.357-1.615.849-.232 1.574-.787 2.132-1.41.56-.627.914-1.28 1.039-1.639.199-.575.356-1.539.428-2.59z"/></svg>'
    return HTMLResponse(content = thumb_up_filled, status_code = 200)

@app.put("/downvote/{event_id}")
async def downvote(event_id: str = None):
    time.sleep(1)
    thumb_down_filled = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-hand-thumbs-down-fill" viewBox="0 0 16 16"><path d="M6.956 14.534c.065.936.952 1.659 1.908 1.42l.261-.065a1.38 1.38 0 0 0 1.012-.965c.22-.816.533-2.512.062-4.51q.205.03.443.051c.713.065 1.669.071 2.516-.211.518-.173.994-.68 1.2-1.272a1.9 1.9 0 0 0-.234-1.734c.058-.118.103-.242.138-.362.077-.27.113-.568.113-.856 0-.29-.036-.586-.113-.857a2 2 0 0 0-.16-.403c.169-.387.107-.82-.003-1.149a3.2 3.2 0 0 0-.488-.9c.054-.153.076-.313.076-.465a1.86 1.86 0 0 0-.253-.912C13.1.757 12.437.28 11.5.28H8c-.605 0-1.07.08-1.466.217a4.8 4.8 0 0 0-.97.485l-.048.029c-.504.308-.999.61-2.068.723C2.682 1.815 2 2.434 2 3.279v4c0 .851.685 1.433 1.357 1.616.849.232 1.574.787 2.132 1.41.56.626.914 1.28 1.039 1.638.199.575.356 1.54.428 2.591"/></svg>'
    return HTMLResponse(content = thumb_down_filled, status_code = 200)

