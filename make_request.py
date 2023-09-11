import requests
import json
import yaml
import argparse

def convert_yaml_to_json(yaml_file):
    try:
        # Read the YAML file
        with open(yaml_file, 'r') as file:
            yaml_data = yaml.load(file, Loader=yaml.FullLoader)

        # Convert YAML to JSON
        json_data = json.dumps(yaml_data, indent=4)

        return json_data
    except FileNotFoundError:
        return None
    

def main():
    parser = argparse.ArgumentParser(description='Convert a YAML file to JSON.')
    parser.add_argument('yaml_file', type=str, help='Path to the YAML file to convert')

    args = parser.parse_args()

    yaml_file_path = args.yaml_file
    json_data = convert_yaml_to_json(yaml_file_path)

    if json_data:
        #print(json_data)
        url = "http://127.0.0.1:12345/spmigration"
        headers = {
        'Content-Type': 'application/json'
        }
        response = requests.request("POST", url, headers=headers, data=json_data)
        print(response.text)
    else:
        print(f"Error: The file '{yaml_file_path}' was not found.")

if __name__ == "__main__":
    main()