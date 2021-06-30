<img src="https://secure.meetupstatic.com/photos/event/5/4/8/a/600_331101642.jpeg" width="150px" align="right" />

# Go Blog Backend

### **Description**

A Go backend for a blog app using JSON Web Token for authentication and mongoDB as database

### **Requirements**

- Go 1.15+
- MongoDB Community Server or MongoDB Atlas account (https://www.mongodb.com/)

### **Install**

```console
$ npm install
```

### **Usage**

- Edit the config file and set your Watson credentials:

```console
$ vi ./config/default.json
```

```json
{
  "server": { "port": "3004" },
  "watson": {
    "apikey": "yourApiKey",
    "serviceUrl": "yourServiceUrl",
    "classifierId": "yourClassifierId"
  }
}
```

_\*Credentials are provided by IBM Watson once the model was trained with the dataset provided in **/datasets/meliDataset.csv**_

- Run App:

```console
$ npm run start
```

- Get the classes of the model used to classify the questions:

```http
GET http://localhost:3004/classes
```

Response:

```javascript
[
  {
    class: "thanks",
    description: "Short answers to a previous reply.",
  },
  {
    class: "stock",
    description: "When the questions is about available stock.",
  },
  {
    class: "payment",
    description: "Questions regarding payment terms.",
  },
  {
    class: "warranty",
    description: "Questions regarding warranty conditions.",
  },
  {
    class: "delivery",
    description: "Questions regarding delivery terms and zones.",
  },
  {
    class: "variation",
    description: "Questions regarding proudct variations.",
  },
  {
    class: "attributes",
    description: "Questions regarding the attributes of the product published.",
  },
  {
    class: "undefined",
    description:
      "Questions that can´t be classified in any of the defined classes",
  },
];
```

- POST a question:

```http
POST http://localhost:3004/question
```

body:

```json
{ "question": "tenes stock?" }
```

Response:

```json
{
  "question": "tenes stock?",
  "classification": [
    {
      "class_name": "stock",
      "confidence": 0.9783917495954597
    },
    {
      "class_name": "variation",
      "confidence": 0.011924305012071207
    },
    {
      "class_name": "warranty",
      "confidence": 0.003217728841590565
    },
    {
      "class_name": "undefined",
      "confidence": 0.0017209188862953393
    },
    {
      "class_name": "thanks",
      "confidence": 0.0015121900697123954
    },
    {
      "class_name": "payment",
      "confidence": 0.0012185533665207027
    },
    {
      "class_name": "attributes",
      "confidence": 0.001185499200915812
    },
    {
      "class_name": "delivery",
      "confidence": 0.0008290550274342253
    }
  ]
}
```
