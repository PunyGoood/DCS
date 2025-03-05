import requests
from urllib.parse import urljoin, urlparse
from bs4 import BeautifulSoup
# flink/src/utils/robots_parser.py

class RobotsParser:
    def __init__(self, base_url):
        self.base_url = base_url
        self.rules = {}
        self._fetch_robots_txt()

    def _fetch_robots_txt(self):
        robots_url = urljoin(self.base_url, '/robots.txt')
        try:
            response = requests.get(robots_url)
            if response.status_code == 200:
                self._parse_robots_txt(response.text)
            else:
                print(f"Failed to fetch robots.txt from {robots_url}")
        except Exception as e:
            print(f"Error fetching robots.txt: {e}")

    def _parse_robots_txt(self, content):
        lines = content.splitlines()
        current_user_agent = None

        for line in lines:
            line = line.strip()
            if not line or line.startswith('#'):
                continue

            parts = line.split(':', 1)
            if len(parts) != 2:
                continue

            key, value = [p.strip() for p in parts]
            if key.lower() == 'user-agent':
                current_user_agent = value
                if current_user_agent not in self.rules:
                    self.rules[current_user_agent] = {'allow': [], 'disallow': []}
            elif key.lower() == 'allow' and current_user_agent:
                self.rules[current_user_agent]['allow'].append(value)
            elif key.lower() == 'disallow' and current_user_agent:
                self.rules[current_user_agent]['disallow'].append(value)

    def is_allowed(self, url):
        parsed_url = urlparse(url)
        path = parsed_url.path

        user_agents = ['*', '*']  # 默认检查所有用户代理
        for user_agent, rules in self.rules.items():
            if user_agent == '*' or user_agent in user_agents:
                for disallow in rules['disallow']:
                    if path.startswith(disallow):
                        return False
                for allow in rules['allow']:
                    if path.startswith(allow):
                        return True
        return True  # 如果没有明确的规则，认为是可以访问的