import express, { Express } from 'express'
import axios from 'axios'
import {
  MeterProvider,
  PeriodicExportingMetricReader,
} from '@opentelemetry/sdk-metrics'
import { OTLPMetricExporter } from '@opentelemetry/exporter-metrics-otlp-http'
import { Resource } from '@opentelemetry/resources'

const morgan = require('morgan')

const metricExporter = new OTLPMetricExporter({
  url: process.env.OTEL_EXPORTER_OTLP_METRICS_ENDPOINT,
})

const meterProvider = new MeterProvider({
  resource: new Resource({
    ['service.name']: process.env.OTEL_SERVICE_NAME || 'node-otlp',
  }),
  readers: [
    new PeriodicExportingMetricReader({
      exporter: metricExporter,
      exportIntervalMillis: 1000,
    }),
  ],
})

const meter = meterProvider.getMeter('custom_metrics')
const requestCounter = meter.createCounter('request.counter', {
  description: 'A request counter',
})

const PORT: number = parseInt(process.env.PORT || '8003')
const app: Express = express()

app.use(morgan('combined'))

app.get('/', (_req, res) => {
  const start = Date.now()

  axios
    .get('https://bible-api.com/?random=verse')
    .then(r => {
      const finish = Date.now()
      const time = (finish - start) / 1000

      requestCounter.add(1, { duration: time })

      res.send(r.data).status(r.status)
    })
    .catch(e => {
      console.error(e)
      res.send('👀 Oops, unknown error ').status(500)
    })
})

app.listen(PORT, () => {
  console.log(`🔥 Listening for requests on http://localhost:${PORT}`)
})
