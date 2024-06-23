import "./app.pcss"
import App from './App.svelte'
import "./lib/services/services"

const app = new App({
  target: document.getElementById('app')!,
})

export default app
