import express, { Express } from 'express'
import axios from 'axios'

const morgan = require('morgan')

const PORT: number = parseInt(process.env.PORT || '8003')
const app: Express = express()

app.use(morgan('combined'))

app.get('/', (_req, res) => {
  axios.get('https://bible-api.com/?random=verse').then(r => {
    res.send(r.data).status(r.status)
  }).catch(e => {
    console.error(e)
    res.send("👀 Oops, unknown error ").status(500)
  })
})

app.listen(PORT, () => {
  console.log(`🔥 Listening for requests on http://localhost:${PORT}`)
})
