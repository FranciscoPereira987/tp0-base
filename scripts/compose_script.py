import jinja2 as jinja
import argparse

COMPOSE_FILE = "docker-compose-dev.yaml"
TEMPLATE_FOLDER = "scripts/"
TEMPLATE_FILE = "docker-compose-template.yaml.jinja"

def create_parser() -> argparse.ArgumentParser:
    """
        Returns an argument parser that looks for either a -c or --clients argument
    """
    parser = argparse.ArgumentParser()

    parser.add_argument("-c", "--clients", type=int, default=2)

    return parser

def make_compose_file(clients: int):
    """
        Loads the template at TEMPLATE_FOLDER/TEMPLATE_FILE, renders it to create a docker-compose 
        file with the desired amount of clients
    """
    loader = jinja.FileSystemLoader(TEMPLATE_FOLDER)
    template = jinja.Environment(loader=loader).get_template(TEMPLATE_FILE)
    with open(COMPOSE_FILE, 'w') as compose_file:
        rendered = template.render(clients=clients)
        compose_file.write(rendered)


if __name__ == "__main__":
    parser = create_parser()
    args = parser.parse_args()
    make_compose_file(args.clients)

