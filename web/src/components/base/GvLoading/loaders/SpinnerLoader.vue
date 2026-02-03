<template>
  <svg class="gv-loader gv-loader--spinner" viewBox="0 0 50 50">
    <circle
      cx="25"
      cy="25"
      r="20"
      fill="none"
      stroke="currentColor"
      stroke-width="4"
      stroke-linecap="round"
      stroke-dasharray="31.4 31.4"
    />
  </svg>
</template>

<style scoped>
.gv-loader--spinner {
  width: 100%;
  height: 100%;
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  100% {
    transform: rotate(360deg);
  }
}
</style>
