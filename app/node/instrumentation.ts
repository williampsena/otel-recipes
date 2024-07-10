import * as opentelemetry from '@opentelemetry/sdk-node'
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node'

const sdk = new opentelemetry.NodeSDK({
  instrumentations: [getNodeAutoInstrumentations()]
})

sdk.start()
