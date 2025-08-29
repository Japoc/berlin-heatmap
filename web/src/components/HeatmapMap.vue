<script setup lang="ts">
import {defineProps, onMounted, onUnmounted, ref, toRefs, watch} from 'vue'
import {LMap, LTileLayer, LImageOverlay, LControlScale, LMarker} from '@vue-leaflet/vue-leaflet'
import SpinnerOverlay from "./SpinnerOverlay.vue";
import { decode } from '@googlemaps/polyline-codec';
import Polyline from "./Polyline.vue";

// store last mouse position on map (for "m" key)
const lastMouseLat = ref<number | null>(null)
const lastMouseLon = ref<number | null>(null)

const metroRoutes = ref<any[]>([])
const sBahnRoutes = ref<any[]>([])
const tramRoutes = ref<any[]>([])
const busRoutes = ref<any[]>([])

// store heatmap url
const heatmapUrl = ref<string | null>(null)
const latitude = ref<string | null>(null)
const longitude = ref<string | null>(null)
const routeData = ref<any | null>(null)

const isLoading = ref(false)
const showSpinner = ref(false)

const apiURL = import.meta.env.VITE_API_URL

const bbox: [[number, number], [number, number]] = [[52.33, 13.08], [52.67, 13.76]]

const props = defineProps<{
  showSBahn: boolean,
  showMetro: boolean,
  showTram: boolean,
  showBus: boolean,
  maxHeatValue: number,
  minHeatValue: number,
}>()

const { showSBahn, showMetro, showTram, showBus, maxHeatValue } = toRefs(props)

async function onMapClick(event: any) {
  longitude.value = event.latlng.lng
  latitude.value = event.latlng.lat
  console.log(longitude.value, latitude.value)
  updateHeatmapUrl()
}

function onMapMouseMove(event: any) {
  // only track mouse position
  lastMouseLat.value = event.latlng.lat
  lastMouseLon.value = event.latlng.lng
}

async function fetchRoutes(routeType: string) {
    const url = `${apiURL}/routes?type=${routeType}`
    const res = await fetch(url)
    if (!res.ok) {
      console.log(new Error(`Failed to fetch ${url}: ${res.status}`))
      return []
    }
    return await res.json()
}

// handle key press ("m" = fetch route from clicked point to mouse position), for debugging purposes
async function onKeyDown(e: KeyboardEvent) {
  if (e.key.toLowerCase() === "m" && latitude.value && longitude.value && lastMouseLat.value && lastMouseLon.value) {
    await fetchRoute(Number(latitude.value), Number(longitude.value), lastMouseLat.value, lastMouseLon.value)
    console.log(routeData.value.distance)
  }
}

let spinnerTimer: number | null = null

function updateHeatmapUrl() {
  if (latitude.value && longitude.value) {
    const url = `${apiURL}/heatmap?lat=${latitude.value}&lon=${longitude.value}&format=png&max=${maxHeatValue.value}`
    heatmapUrl.value = url
    isLoading.value = true

    // Show spinner only if request > 1s
    if (spinnerTimer) clearTimeout(spinnerTimer)
    spinnerTimer = window.setTimeout(() => {
      if (isLoading.value) showSpinner.value = true
    }, 1000)

    const img = new Image()
    img.src = url
    img.onload = () => {
      heatmapUrl.value = url
      isLoading.value = false
      showSpinner.value = false
      if (spinnerTimer) clearTimeout(spinnerTimer)
    }
    img.onerror = () => {
      isLoading.value = false
      showSpinner.value = false
      if (spinnerTimer) clearTimeout(spinnerTimer)
    }

  }


}

watch(maxHeatValue, () => {
  updateHeatmapUrl()
})

async function fetchRoute(latFrom: number, lonFrom: number, latTo: number, lonTo: number) {
  console.log(latFrom, lonFrom, latTo, lonTo)
  const routeUrl = `${apiURL}/route?latFrom=${latFrom}&lonFrom=${lonFrom}&latTo=${latTo}&lonTo=${lonTo}`
  try {
    const res = await fetch(routeUrl)
    if (!res.ok) throw new Error(`Route fetch failed: ${res.status}`)
    routeData.value = await res.json()
  } catch (err) {
    console.error(err)
    routeData.value = { error: String(err) }
  }
}

onMounted(async () => {
  window.addEventListener("keydown", onKeyDown)
  metroRoutes.value = await fetchRoutes("metro")
  sBahnRoutes.value = await fetchRoutes("sbahn")
  tramRoutes.value = await fetchRoutes("tram")
  busRoutes.value = await fetchRoutes("bus")
})
onUnmounted(() => {
  window.removeEventListener("keydown", onKeyDown)
  //<iframe src="https://japoc.github.io/berlin-heatmap" width="500px" height="500px"></iframe>
})

</script>

<template>
    <div class="map-container">
      <LMap
          :zoom="11"
          :center="[52.52, 13.40]"
          @click="onMapClick"
          @mousemove="onMapMouseMove"
      >
        <LMarker v-if="latitude && longitude" :lat-lng="[Number(latitude), Number(longitude)]" />
        <LTileLayer
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            attribution="&copy; OpenStreetMap contributors"
        />
        <Polyline
            v-if="showMetro"
            v-for="(item) in metroRoutes"
            :mode="item.Mode"
            :short-name="item.ShortName"
            :color="item.Color"
            :points="decode(item.Points)" >
        </Polyline>
        <Polyline
            v-if="showSBahn"
            v-for="(item) in sBahnRoutes"
            :mode="item.Mode"
            :short-name="item.ShortName"
            :color="item.Color"
            :points="decode(item.Points)" >
        </Polyline>
        <Polyline
            v-if="showTram"
            v-for="(item) in tramRoutes"
            :mode="item.Mode"
            :short-name="item.ShortName"
            :color="'D70040'"
            :points="decode(item.Points)" >
        </Polyline>
        <Polyline
            v-if="showBus"
            v-for="(item) in busRoutes"
            :mode="item.Mode"
            :short-name="item.ShortName"
            :color="'9F2B68'"
            :points="decode(item.Points)" >
        </Polyline>
        <LImageOverlay
            v-if="heatmapUrl"
            :url="heatmapUrl"
            :bounds="bbox"
            :opacity="0.5"
        />
        <LControlScale :imperial="false" />
      </LMap>
      <SpinnerOverlay v-if="showSpinner" text="Waking up backend... please wait" />
    </div>
</template>


<style scoped>

.map-container {
  flex: 1;
  position: relative;
}

.map-container .leaflet-container {
  height: 100%;
  width: 100%;
}


@keyframes spin {
  to { transform: rotate(360deg); }
}

:global(.leaflet-interactive) {
  outline: none !important;
  stroke-opacity: 1;    /* ensure opacity stays correct */
}
</style>
