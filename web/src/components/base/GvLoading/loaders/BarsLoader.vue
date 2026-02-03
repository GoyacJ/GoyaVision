<template>
  <div class="gv-loader gv-loader--bars">
    <span class="gv-loader__bar"></span>
    <span class="gv-loader__bar"></span>
    <span class="gv-loader__bar"></span>
    <span class="gv-loader__bar"></span>
  </div>
</template>

<style scoped>
.gv-loader--bars {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 100%;
}

.gv-loader__bar {
  width: 15%;
  height: 100%;
  background-color: currentColor;
  border-radius: 4px;
  animation: stretch 1.2s infinite ease-in-out;
}

.gv-loader__bar:nth-child(1) {
  animation-delay: -0.48s;
}

.gv-loader__bar:nth-child(2) {
  animation-delay: -0.32s;
}

.gv-loader__bar:nth-child(3) {
  animation-delay: -0.16s;
}

@keyframes stretch {
  0%, 40%, 100% {
    transform: scaleY(0.4);
    opacity: 0.5;
  }
  20% {
    transform: scaleY(1);
    opacity: 1;
  }
}
</style>
