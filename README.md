**HOW TO USE THIS PROJECT**

*INSTRUCTIONS FOR RUNNING ON UBUNTU:*

*Note: It is assumed that Golang is already installed on your computer*
1. Install curl if you don't have
```
sudo apt install curl
```
2. Download and install Elasticsearch 
```
go get github.com/elastic/go-elasticsearch/v9@latest
```
3. Launch and check ES
```
systemctl start elasticsearch.service
systemctl status elasticsearch.service
```
4. Make a simple request for make sure
```
curl http://localhost:9200
```
You should get something like this:
```
{
  "name" : "Your machine name",
  "cluster_name" : "elasticsearch",
  "cluster_uuid" : "kshdkfhslkf",
  "version" : {
    "number" : "9.1.5",
    "build_flavor" : "default",
    "build_type" : "deb",
    "build_hash" : "09843kljslk209840092jf0",
    "build_date" : "2025-10-02T22:07:12.966975992Z",
    "build_snapshot" : false,
    "lucene_version" : "10.2.2",
    "minimum_wire_compatibility_version" : "8.19.0",
    "minimum_index_compatibility_version" : "8.0.0"
  },
  "tagline" : "You Know, for Search"
}
```

**NOTE: before next steps your need to set an evironment variable with your secret key for autentify with JWT**

Linux/macOS (bash): 
```
export ES_API_KEY="your-secret-key"
```
Windows (PowerShell):
```
$env:ES_API_KEY = "your-secret-key"
```

5. Now you need to make indices and mapping for ES, use next command from src folder:
```
make indexcreator
// you can write your own index by launching utility with command-line, usage: $utility -i idxname
```

NOTE: you can use next command for default data and configs, if you want just to check how it works:
```
make
```

**If you need you own custom config and data just follow to next steps.**

6. You a ready to download a data to the database
```
make dataloader
// you can write your own path/to/file.csv to load and choose index by launching utility with command-line
// usage: $utility -i index_to_load -f path/to/file.csv
```

7. You can to handle Elsaticsearch pagination for more than 20 000 entires (config your own number) with next command:
```
~$ curl -XPUT -H "Content-Type: application/json" "http://localhost:9200/places/_settings" -d '
{
  "index" : {
    "max_result_window" : 20000 // write here your number
  }
}'
// NOTE: you need to replace 'places' with your indexname if you used a custom
```

8. Let's check our data in browser. Make next command and write in browser's search-line http://localhost:8888
```
make http
```

9. OK, now let's finish with autentification with JWT. Write in browser's search-line http://localhost:8888/api/get_token.
Copy token and use curl or Postman to check the JWT Autentification is working, example:
```
curl -H "Authorization: Bearer <your_token>" http://localhost:8888/api/recommend?lat=55.674&lon=37.666
```


**Task in russian:**

1. Загрузка данных. На рынке существует множество разных баз данных. Но поскольку мы пытаемся обеспечить возможность поиска объектов, давайте использовать Elasticsearch. Elasticsearch — полнотекстовый движок поиска, построенный на Lucene. Он предоставляет HTTP API, которое мы будем использовать в этом задании. Наш предоставленный набор данных о ресторанах (взятый из портала открытых данных) состоит из более чем 13 тысяч ресторанов в Москве, Россия (вы можете сформировать другой похожий набор данных для любого другого места, которое захотите, и передавать его в модуль **dataloader** в качестве параметра с флагом -f). Каждая запись содержит:

- ID
- Название
- Адрес
- Телефон
- Долгота
- Широта

Прежде чем загружать все записи в базу данных, давайте создадим индекс и сопоставление (явно указывая типы данных). Без них Elasticsearch будет пытаться угадать типы полей по предоставленным документам, и иногда он может не распознать геопозицию. schema.json:
```
{
  "properties": {
    "name": {
        "type":  "text"
    },
    "address": {
        "type":  "text"
    },
    "phone": {
        "type":  "text"
    },
    "location": {
      "type": "geo_point"
    }
  }
}
```

В этом задании вы должны использовать привязки Go к Elasticsearch, чтобы выполнить то же самое.

Далее вам нужно определить отображения типов для наших данных. Это действие должен выполнить программа на Go, которую вы напишете. Теперь у вас есть набор данных для загрузки. Нужно использовать Bulk API для загрузки.

2. Простой интерфейс. Теперь давайте создадим HTML-интерфейс для нашей базы данных. Нам просто нужно отобразить страницу со списком имен, адресов и телефонов, чтобы пользователи могли видеть ее в браузере. Вы должны абстрагировать свою базу данных от HTML-интерфейса:
```
type Store interface {
    // Returns a list of items, a total number of hits and (or) an error in case of one
    GetPlaces(limit int, offset int) ([]types.Place, int, error)
}
```

В основном пакете не должно быть импорта, связанного с Elasticsearch, так как все, что связано с базой данных, должно находиться в пакете db в вашем проекте, и вы должны использовать только этот интерфейс, описанный выше, для взаимодействия с ним. Ваше HTTP-приложение должно запускаться на порту 8888, отвечать списком ресторанов и предоставлять простую разбивку по страницам.

Ссылка "Предыдущая" должна исчезнуть на странице 1, а ссылка "Следующая" - на последней странице.
ВАЖНОЕ ПРИМЕЧАНИЕ: Вы можете заметить, что по умолчанию Elasticsearch не позволяет выполнять разбивку на страницы для более чем 10000 записей.

Кроме того, если параметр "страница" указан с неправильным значением (за пределами [0..last_page] или не числовым), ваша страница должна возвращать сообщение об ошибке HTTP 400 и обычный текст с описанием ошибки:
```
Invalid 'page' value: 'foo'.
```

3. API. В современном мире большинство приложений предпочитают API обычному HTML. Итак, в этом упражнении все, что вам нужно сделать, это реализовать другой обработчик, который отвечает с помощью Content-Type: application/json и JSON-версии того же, что и в п.2. Кроме того, если параметр "страница" указан с неправильным значением (за пределами [0..last_page] или не числовым), ваш API должен выдать соответствующую ошибку.


4. Ближайшие рестораны. Теперь необходимо реализовать главную функцию — поиск **трех** ближайших ресторанов! Для этого вам нужно настроить сортировку по вашему запросу:
```
"sort": [
    {
      "_geo_distance": {
        "location": {
          "lat": 55.674,
          "lon": 37.666
        },
        "order": "asc",
        "unit": "km",
        "mode": "min",
        "distance_type": "arc",
        "ignore_unmapped": true
      }
    }
]
```
где "широта" и "долгота" - это ваши текущие координаты. Итак, для URL-адреса, подобного http://127.0.0.1:8888/api/recommend?lat=55.674&lon=37.666, ваше приложение должно возвращать JSON следующим образом:
```
{
  "name": "Recommendation",
  "places": [
    {
      "id": 30,
      "name": "Ryba i mjaso na ugljah",
      "address": "gorod Moskva, prospekt Andropova, dom 35A",
      "phone": "(499) 612-82-69",
      "location": {
        "lat": 55.67396575768212,
        "lon": 37.66626689310591
      }
    },
    {
      "id": 3348,
      "name": "Pizzamento",
      "address": "gorod Moskva, prospekt Andropova, dom 37",
      "phone": "(499) 612-33-88",
      "location": {
        "lat": 55.673075576456,
        "lon": 37.664533747576
      }
    },
    {
      "id": 3347,
      "name": "KOFEJNJa «KAPUChINOFF»",
      "address": "gorod Moskva, prospekt Andropova, dom 37",
      "phone": "(499) 612-33-88",
      "location": {
        "lat": 55.672865251005106,
        "lon": 37.6645689561318
      }
    }
  ]
}
```

5. JWT. Последнее, что нужно сделать, - это обеспечить простую форму аутентификации. В настоящее время одним из самых популярных способов реализации этого для API является использование JWT. К счастью, в Go есть довольно хороший набор инструментов для решения этой задачи.
Во-первых, вам нужно реализовать конечную точку API http://127.0.0.1:8888/api/get_token единственной целью которой будет сгенерировать токен и вернуть его следующим образом:
```
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNjAxOTc1ODI5LCJuYW1lIjoiTmlrb2xheSJ9"
}
```
Не забудьте указать заголовок "Content-Type: application/json".
Во-вторых, вам необходимо защитить свою конечную точку /api/recommend с помощью промежуточного программного обеспечения JWT, которое проверяет действительность этого токена.
Таким образом, по умолчанию, когда этот API запрашивается из браузера, он должен завершаться ошибкой HTTP 401, но работать, если клиент предоставляет заголовок Authorization: Bearer <токен> (вы можете проверить это с помощью cURL или Postman).
Это самый простой способ обеспечить аутентификацию.