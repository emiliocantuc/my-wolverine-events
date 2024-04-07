import json


if __name__ == '__main__':

    with open('static/json/programs_grad.json', 'r') as f:
        programs_grad = json.load(f)
    
    with open('static/json/programs_undergrad.json', 'r') as f:
        programs_undergrad = json.load(f)

    with open('templates/preferences.html', 'r') as f:
        template = f.read()

    programs_HTML = []
    for program in programs_grad + programs_undergrad:
        programs_HTML.append(f'<option value="{program}">{program}</option>')

    template = template.replace('{{programs}}', '\n\t\t\t'.join(programs_HTML))
    
    with open('preferences.html', 'w+') as f:
        f.write(template)
