import requests, os, json
from bs4 import BeautifulSoup


def parse_undergrad_programs(html):

    program_names = []
    soup = BeautifulSoup(html, 'html.parser')

    # Get first div with class = "view-content"
    div = soup.find('div', class_='view-content')

    # Iterate over all tables in div
    for table in div.find_all('table'):

        # and iterate over all rows in table
        for row in table.find_all('tr'):
            cols = row.find_all('td')
            if len(cols) < 2:
                continue
            major, school = cols
            major, school = major.text.strip(), school.text.strip()
            program_names.append(major)

    return program_names


def get_undergrad_programs(url):
    r = requests.get(url)
    if r.status_code != 200:
        return None
    return parse_undergrad_programs(r.text)


def parse_grad_programs(html):

    program_names = []
    soup = BeautifulSoup(html, 'html.parser')

    # get table with attribute aria-label = "Programs of Study"
    table = soup.find('table', {'aria-label': 'Programs of Study'})

    # iterate through table rows
    for row in table.find_all('tr'):
        
        cells = row.find_all('td')
        # if more than one cell
        if len(cells) > 1:
            name, campus, school, deg_type = cells[:4]
            name, campus, school, deg_type = name.text.strip(), campus.text.strip(), school.text.strip(), deg_type.text.strip()
            program_names.append(name)
    
    return program_names

def get_grad_programs(url):
    r = requests.get(url)
    if r.status_code != 200:
        return None
    return parse_grad_programs(r.text)

def write_programs(programs, filename):
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    with open(filename, 'w+') as f:
        f.write(json.dumps(programs, indent = 4))
    print(f'Wrote {len(programs)} programs to {filename}')



if __name__ == '__main__':

    UNDERGRAD_PROGRAMS_URL = 'https://admissions.umich.edu/academics-majors/majors-degrees'
    GRAD_PROGRAMS_URL = 'https://rackham.umich.edu/programs-of-study/'
    JSON_DIR = 'static/json/'

    # Get undergrad programs
    undergrad_programs = get_undergrad_programs(UNDERGRAD_PROGRAMS_URL)
    if undergrad_programs:
        write_programs(undergrad_programs, os.path.join(JSON_DIR,'programs_undergrad.json'))
    else:
        print('Failed to get undergrad programs')

    # Get grad programs
    grad_programs = get_grad_programs(GRAD_PROGRAMS_URL)
    if grad_programs:
        write_programs(grad_programs, os.path.join(JSON_DIR,'programs_grad.json'))
    else:
        print('Failed to get grad programs')