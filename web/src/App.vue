<script setup lang="ts">
import HeatmapMap from './components/HeatmapMap.vue'
import Slider from "./components/Slider.vue";
import Checkbox from "./components/Checkbox.vue";
import {ref} from "vue";

const minHeatValue = 0
const maxHeatValue = ref(120)

const showMetro = ref(false)
const showSBahn = ref(false)
</script>

<template>

  <div class="page-container" style="display: flex; flex-direction: column; height: 100vh;">
    <HeatmapMap
        :max-heat-value="maxHeatValue"
        :min-heat-value="minHeatValue"
        :show-metro="showMetro"
        :show-s-bahn="showSBahn"
    >
    </HeatmapMap>
    <div class="bottom-panel">
      <div class="bottom-controls">

        <!-- Left block -->
        <div class="left-block">
          <div class="heatmap-legend-below">
            <div class="color-bar"></div>
            <div class="labels">
              <span>{{ minHeatValue }}</span>
              <span>{{ Math.round((minHeatValue + maxHeatValue)/2) }}</span>
              <span>{{ maxHeatValue }}</span>
            </div>
          </div>
          <Slider v-model="maxHeatValue" :min="10" :max="180" :step="10" :label="'maximum heat'"/>
        </div>

        <!-- Right block -->
        <div class="right-block">
          <Checkbox v-model="showMetro" label="Show Metro Routes" />
          <Checkbox v-model="showSBahn" label="Show S-Bahn Routes" />
        </div>

      </div>
    </div>
  </div>
</template>

<style>
.page-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}

.bottom-panel {
  padding: 10px;
  background: #5e5e5e;
}

.bottom-controls {
  display: flex;
  justify-content: space-between;
  align-items: flex-start; /* align top edges */
  gap: 20px;
}

/* Left block: legend stacked above slider */
.left-block {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  min-width: 250px;
}

/* Right block: checkboxes stacked vertically */
.right-block {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 180px;
}

.heatmap-legend-below .color-bar {
  height: 15px;
  width: 100%;
  background: linear-gradient(to right, blue, cyan, yellow, red);
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

