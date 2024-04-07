import json


if __name__ == '__main__':

    # load json
    with open('static/json/json.json', 'r') as f:
        events = json.load(f)

    with open('templates/card.html', 'r') as f:
        card_template = f.read()

    # cut off text to l chars and add '...'
    cutoff = lambda s, l: s if len(s) <= l else s[:l] + '...'
    
    cards_html = []
    for event in events:
        card = card_template[::]
        card = card.replace('{{title}}', cutoff(event['event_title'], 50))
        card = card.replace('{{description}}', cutoff(event['description'], 200))
        cards_html.append(card)

    with open('templates/index.html', 'r') as f:
        index = f.read()
        index = index.replace('{{events}}', '\n\t\t\t'.join(cards_html))
    
    with open('index.html', 'w+') as f:
        f.write(index)
