import requests
from bs4 import BeautifulSoup
import time
from random import choice
from fake_useragent import UserAgent
from urllib.request import Request, urlopen
import urllib

ua = UserAgent()
proxies = []

def refresh_proxies():
   # proxies_req = Request('https://www.sslproxies.org/')
   # proxies_req.add_header('User-Agent', ua.random)
   # proxies_doc = urlopen(proxies_req).read().decode('utf8')
    proxies_doc = requests.get('https://www.sslproxies.org/', headers=get_random_user_agent()).content
    soup = BeautifulSoup(proxies_doc, 'html.parser')
    proxies_table = soup.find(id='proxylisttable')

    # Save proxies in the array
    for row in proxies_table.tbody.find_all('tr'):
        proxies.append({
            'ip':   row.find_all('td')[0].string,
            'port': row.find_all('td')[1].string
        })
    
def get_random_proxy():
    if not proxies:
        refresh_proxies()
    
    proxy_raw = choice(proxies)
    proxy = {"http": "{}:{}".format(proxy_raw['ip'], proxy_raw['port'])}
    #proxy = '{}:{}'.format(proxy_raw['ip'], proxy_raw['port'])
    return proxy

def get_random_user_agent():
    user_agent = {"User-Agent": ua.random}
    return user_agent

def fetch_results(search_term, number_results, language_code):
    assert isinstance(search_term, str), 'Search term must be a string'
    assert isinstance(number_results, int), 'Number of results must be an integer'
    escaped_search_term = search_term.replace(' ', '+')

    google_url = 'https://www.google.com/search?q={}&num={}&hl={}'.format(escaped_search_term, number_results, language_code)
    response = requests.get(google_url, headers=get_random_user_agent(), proxies=get_random_proxy())
    response.raise_for_status()
   # req = Request(google_url)
   # req.add_header('User-Agent', ua.random)

    #proxy_raw = get_random_proxy()
    #proxy = proxy_raw['ip'] + ':' + proxy_raw['port']
    #req.set_proxy(get_random_proxy(), 'http')
   # try:
        #response = urlopen(req).read().decode('utf8')
   # except urllib.error.HTTPError as e:
   ##     print(e)
    #except urllib.error.URLError as e:
    #    print(e)
    #except Exception:
    #    import traceback
    #    print(traceback.format_exc())

    return search_term, response

def parse_results(html, keyword):
    soup = BeautifulSoup(html, 'html.parser')

    found_results = []
    rank = 1
    result_block = soup.find_all('div', attrs={'class': 'g'})
    for result in result_block:

        link = result.find('a', href=True)
        title = result.find('h3', attrs={'class': 'r'})
        description = result.find('span', attrs={'class': 'st'})
        if link and title:
            link = link['href']
            title = title.get_text()
            if description:
                description = description.get_text()
            if link != '#':
                found_results.append({'keyword': keyword, 'rank': rank, 'title': title, 'description': description, 'link': link})
                rank += 1
    return found_results


def scrape_google(search_term, number_results, language_code):
    try:
        keyword, html = fetch_results(search_term, number_results, language_code)
        results = parse_results(html, keyword)
        return results
    except AssertionError:
        raise Exception("Incorrect arguments parsed to function")
    except requests.HTTPError:
        raise Exception("You appear to have been blocked by Google")
    except requests.RequestException:
        raise Exception("Appears to be an issue with your connection")


if __name__ == '__main__':
   # for keyword in keywords:
    data = []
    keyword = "site%3Aspotify.com+inurl%3Aplaylist+voxtrot+nujabes"
    try:
        results = scrape_google(keyword, 100, "en")
        for result in results:
            data.append(result)
    except Exception as e:
        print(e)
    #finally:
       # time.sleep(10)
    print(data)