<template>
  <div class="player-container">
      <div class="flex">
        <div class="left-box">
          <span v-if="isAlive">{{ player.state.health }}</span>
          <img v-else src="/csgo-icons/skull-solid.svg" class="svg-icon"/>
        </div>
        <div class="right-box">
          <div class="health" :style="`background: linear-gradient(to left, black ${100-player.state.health}%, ${player.team.side === 'T' ? '#f52a2a' : '#3b5fed'} 0%);`">
            <div>{{ player.name }}</div>
          </div>
          <div>
            <div>
              <img src="/csgo-icons/crosshairs-solid.svg" class="svg-icon"/>
              {{player.stats.kills}}
              <img src="/csgo-icons/skull-solid.svg" class="svg-icon"/>
              {{player.stats.deaths}}
            </div>
            <span v-if="isAlive">
              <span v-for="(weapon, key) in player.weapons" :key="weapon">
                <img
                  v-if="key === 'weapon_0'"
                  :src="`/csgo-icons/${(weapon.state === 'active' ? 'weapon_knife' : 'weapon_knife_holstered')}.png`"
                  class="weapon-image"
                />
                <img
                  v-else-if="key === 'weapon_1'"
                  :src="`/csgo-icons/${(weapon.state === 'active' ? 'weapon_deagle' : 'weapon_deagle_holstered')}.png`"
                  class="weapon-image"
                />
                <img
                  v-else
                  :src="`/csgo-icons/${(weapon.state === 'active' ? weapon.name : weapon.name+'_holstered')}.png`"
                  class="weapon-image"
                />
              </span>

            </span>
          </div>
        </div>
      </div>
  </div>
</template>


<style scoped>
* {
  color: white;
}
.right-box {
  display: flex;
  flex-direction: column;
  background-color: black;
}
.left-box {
  background-color: black;
  width: 50px;
  display: flex;
  justify-content: center;
  align-items: center;
}
.flex {
  display: flex;
}
.health {
  padding: .5rem;
  width: 250px;
}
.player-container {
  padding: 1rem
}
.weapon-image {
  width: 25x;
  height: 25x;
}
.svg-icon {
  filter: invert(100%) sepia(5%) saturate(2878%) hue-rotate(0deg) brightness(100%) contrast(100%);
  height: 20px;
  width: 20px;
}
</style>

<script>
export default {
  props: ['player'],
  data() {
    return {}
  },
  mounted() {  },
  computed: {
    isAlive() {
      return this.player.state.health > 0
    },
  }
}
</script>