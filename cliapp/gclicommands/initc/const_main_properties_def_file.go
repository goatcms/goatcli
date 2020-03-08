package initc

// PropertiesDefFile is content of .goat/properties.def/00_main.json file
const PropertiesDefFile = `[{
  "prompt":"Insert application userfriendly name",
  "key":"app.name",
  "type":"line",
  "min":1,
  "max":150
},{
  "prompt":"Insert application slogan",
  "key":"app.slogan",
  "type":"line",
  "min":1,
  "max":500
},{
  "prompt":"Insert application description",
  "key":"app.description",
  "type":"line",
  "min":1,
  "max":5000
}]`
