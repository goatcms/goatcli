package initc

// AuthorDataDefFile is content of .goat/data.def/appauthor.json file
const AuthorDataDefFile = `[{
  "type":"appauthor",
  "properties":[{
    "prompt":"Firstname",
    "key":"firstname",
    "type":"line",
    "min":1,
    "max":500
  }, {
    "prompt":"Lastname",
    "key":"lastname",
    "type":"line",
    "min":1,
    "max":500
  }, {
    "prompt":"Email",
    "key":"email",
    "type":"line",
    "pattern": "^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$",
    "min":1,
    "max":600
  }]
}]`
