<script setup lang="ts">
import {LPolyline} from "@vue-leaflet/vue-leaflet";
import type {LatLngExpression, LeafletMouseEvent} from "leaflet";

function showTooltip(event: LeafletMouseEvent, text: string) {
  const polyline = event.target; // Get the polyline that triggered the event
  polyline.bindTooltip(text, {
    permanent: false, // Tooltip only shows on hover
    direction: 'top', // Tooltip position relative to the polyline
    sticky: true,
    offset: [0, -10], // Optional offset for better positioning
  }).openTooltip();
}
function hideTooltip(event: LeafletMouseEvent) {
  const polyline = event.target;
  polyline.closeTooltip();
}

defineProps<{
  mode: string,
  shortName: string,
  color: string,
  points: LatLngExpression[],
}>()

</script>

<template>
  <!-- acutal shown line-->
  <LPolyline
      :opacity="0.7"
      :lat-lngs="points"
      :color="'#' + color"
      :name="shortName"
      :weight="3"
      :interactive="false"
  />
  <!-- transparent line with bigger hitbox and hover tooltip-->
  <LPolyline
      :lat-lngs="points"
      color="transparent"
      :name="shortName"
      @mouseover="showTooltip($event, shortName)"
      @mouseout="hideTooltip($event)"
      :weight="12"
      :interactive="true"
  />
</template>
