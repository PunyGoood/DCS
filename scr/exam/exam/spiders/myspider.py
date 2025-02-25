import scrapy

class MySpider(scrapy.Spider):
    name = "myspider"  # 爬虫名称


    start_urls = [
        'https://example.com/page1',
        'https://example.com/page2',
        'https://example.com/page3',
    ]

    def parse(self, response):
        # 提取页面标题
        title = response.css('title::text').get()
        
        # 提取页面中的所有链接
        links = response.css('a::attr(href)').getall()

        # 返回提取的数据
        yield {
            'url': response.url,  # 当前页面的 URL
            'title': title,       # 页面标题
            'links': links,       # 页面中的所有链接
        }

        