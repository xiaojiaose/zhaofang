# Index模版设置

###  Logistics -

    - template Name
       house-resource

    - index Patterns
      house-*

###  index setting -
``` json
{
    "analysis": {
        "analyzer": {
            "my_analyzer": {
                "tokenizer": "whitespace",
                "token_filter": [
                    "lowercase"
                ],
                "type": "custom"
            }
        }
    }
}
 
```

###   mapping -
```json
{
    "properties": {
        "@timestamp": {
            "type": "date",
            "index": true,
            "store": false,
            "sortable": true,
            "aggregatable": true,
            "highlightable": false
        },
        "_id": {
            "type": "keyword",
            "index": true,
            "store": false,
            "sortable": true,
            "aggregatable": true,
            "highlightable": false
        },
        "city": {
            "type": "text",
            "index": true,
            "store": false,
            "sortable": false,
            "aggregatable": true,
            "highlightable": false
        },
        "district_ids": {
            "type": "text",
            "index": true,
            "store": false,
            "sortable": false,
            "aggregatable": true,
            "highlightable": false
        },
        "districts": {
            "type": "text",
            "analyzer": "my_analyzer",
            "index": true,
            "store": false,
            "sortable": false,
            "aggregatable": true,
            "highlightable": false
        },
        "feature": {
            "type": "keyword",
            "index": true,
            "store": false,
            "sortable": false,
            "aggregatable": true,
            "highlightable": false
        },
        "house_id": {
            "type": "text",
            "index": true,
            "store": false,
            "sortable": false,
            "aggregatable": false,
            "highlightable": false
        },
        "house_type": {
            "type": "text",
            "analyzer": "my_analyzer",
            "index": true,
            "store": false,
            "sortable": false,
            "aggregatable": true,
            "highlightable": false
        },
        "last_update": {
            "type": "date",
            "format": "2006-01-02T15:04:05Z07:00",
            "index": true,
            "store": false,
            "sortable": true,
            "aggregatable": true,
            "highlightable": false
        },
        "price": {
            "type": "numeric",
            "index": true,
            "store": false,
            "sortable": true,
            "aggregatable": true,
            "highlightable": false
        },
        "rent_type": {
            "type": "text",
            "analyzer": "my_analyzer",
            "index": true,
            "store": false,
            "sortable": false,
            "aggregatable": true,
            "highlightable": false
        },
        "status": {
            "type": "text",
            "index": true,
            "store": false,
            "sortable": false,
            "aggregatable": true,
            "highlightable": false
        },
        "xiaoqu": {
            "type": "text",
            "index": true,
            "store": false,
            "sortable": false,
            "aggregatable": false,
            "highlightable": false
        },
        "xiaoqu_id": {
            "type": "numeric",
            "index": true,
            "store": false,
            "sortable": true,
            "aggregatable": true,
            "highlightable": false
        }
    }
}
```