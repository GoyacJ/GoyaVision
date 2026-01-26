<template>
  <div class="hls-preview">
    <video
      ref="videoRef"
      class="video-js vjs-default-skin"
      controls
      preload="auto"
      :width="width"
      :height="height"
    ></video>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import videojs from 'video.js'
import 'video.js/dist/video-js.css'

interface Props {
  hlsUrl: string
  width?: number
  height?: number
}

const props = withDefaults(defineProps<Props>(), {
  width: 640,
  height: 360
})

const videoRef = ref<HTMLVideoElement | null>(null)
let player: any = null

watch(() => props.hlsUrl, (newUrl) => {
  if (player && newUrl) {
    player.src({
      src: newUrl,
      type: 'application/x-mpegURL'
    })
    player.play()
  }
})

onMounted(() => {
  if (videoRef.value) {
    player = videojs(videoRef.value, {
      fluid: false,
      responsive: false,
      html5: {
        hls: {
          withCredentials: false
        }
      }
    })

    if (props.hlsUrl) {
      player.src({
        src: props.hlsUrl,
        type: 'application/x-mpegURL'
      })
    }
  }
})

onUnmounted(() => {
  if (player) {
    player.dispose()
  }
})
</script>

<style scoped>
.hls-preview {
  display: inline-block;
}
</style>
