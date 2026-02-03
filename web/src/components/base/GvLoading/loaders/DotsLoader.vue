<template>
  <div class="gv-loader gv-loader--dots">
    <span class="gv-loader__dot"></span>
    <span class="gv-loader__dot"></span>
    <span class="gv-loader__dot"></span>
  </div>
</template>

<style scoped>
.gv-loader--dots {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.gv-loader__dot {
  width: 25%;
  aspect-ratio: 1;
  border-radius: 50%;
  background-color: currentColor;
  animation: bounce 1.4s infinite ease-in-out both;
}

.gv-loader__dot:nth-child(1) {
  animation-delay: -0.32s;
}

.gv-loader__dot:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes bounce {
  0%, 80%, 100% {
    transform: scale(0);
    opacity: 0.5;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}
</style>
