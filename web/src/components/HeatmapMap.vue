<script setup lang="ts">
import {onMounted, onUnmounted, ref, watch} from 'vue'
import {LMap, LTileLayer, LImageOverlay, LControlScale, LMarker} from '@vue-leaflet/vue-leaflet'
import Slider from "./Slider.vue";
// store last mouse position on map (for "m" key)
const lastMouseLat = ref<number | null>(null)
const lastMouseLon = ref<number | null>(null)

// store heatmap url
const heatmapUrl = ref<string | null>(null)
const latitude = ref<string | null>(null)
const longitude = ref<string | null>(null)
const routeData = ref<any | null>(null)

const minHeatValue = 0
const maxHeatValue = ref(120)

const bbox: [[number, number], [number, number]] = [[52.33, 13.08], [52.67, 13.76]]


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

// handle key press ("m" = fetch route from clicked point to mouse position), for debugging purposes
async function onKeyDown(e: KeyboardEvent) {
  if (e.key.toLowerCase() === "m" && latitude.value && longitude.value && lastMouseLat.value && lastMouseLon.value) {
    await fetchRoute(Number(latitude.value), Number(longitude.value), lastMouseLat.value, lastMouseLon.value)
    console.log(routeData.value.distance)
  }
}

function updateHeatmapUrl() {
  if (latitude.value && longitude.value) {
    heatmapUrl.value = `http://localhost:8088/heatmap?lat=${latitude.value}&lon=${longitude.value}&format=png&max=${maxHeatValue.value}`
  }
}

// 🔹 Watch slider value, refresh heatmap when it changes
watch(maxHeatValue, () => {
  updateHeatmapUrl()
})

async function fetchRoute(latFrom: number, lonFrom: number, latTo: number, lonTo: number) {
  console.log(latFrom, lonFrom, latTo, lonTo)
  const routeUrl = `http://localhost:8088/route?latFrom=${latFrom}&lonFrom=${lonFrom}&latTo=${latTo}&lonTo=${lonTo}`
  try {
    const res = await fetch(routeUrl)
    if (!res.ok) throw new Error(`Route fetch failed: ${res.status}`)
    routeData.value = await res.json()
  } catch (err) {
    console.error(err)
    routeData.value = { error: String(err) }
  }
}

onMounted(() => {
  window.addEventListener("keydown", onKeyDown)
})
onUnmounted(() => {
  window.removeEventListener("keydown", onKeyDown)
})
</script>

<template>
  <LMap
      style="height: 100vh; width: 100vh"
      :zoom="11"
      :center="[52.52, 13.40]"
      @click="onMapClick"
      @mousemove="onMapMouseMove"
  >
    <LMarker
        :lat-lng="[Number(latitude), Number(longitude)]"
    ></LMarker>
    <LTileLayer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        attribution="&copy; OpenStreetMap contributors"
    />
    <LImageOverlay
        v-if="heatmapUrl"
        :url="heatmapUrl"
        :bounds="bbox"
        :opacity="0.5"
    />
    <div><a>test</a></div>
    <LControlScale
        :imperial = "false">
    </LControlScale>
  </LMap>
  <div class="heatmap-legend-below">
    <div class="color-bar"></div>
    <div class="labels">
      <span>{{ minHeatValue }}</span>
      <span>{{ Math.round((minHeatValue + maxHeatValue)/2) }}</span>
      <span>{{ maxHeatValue }}</span>
    </div>
  </div>
  <Slider v-model="maxHeatValue" :min="10" :max="180" :step="10" :label="'maximum heat'"/>
</template>

<style scoped>
.heatmap-legend-below {
  width: 100%;
  max-width: 600px;
  margin: 10px auto;
  text-align: center;
  font-size: 12px;
}

.heatmap-legend-below .color-bar {
  height: 15px;
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
  background: linear-gradient(to right, blue, cyan, lime, yellow, red);
  border: 1px solid #aaa;
  border-radius: 2px;
}

.heatmap-legend-below .labels {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  margin-top: 2px;
}

</style>
