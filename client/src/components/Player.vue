<template>
  <div class="player-container">
      <div class="flex" :class="{bordered: focused, 'hidden-border': !focused}">
        <div class="left-box">
          <div v-if="isAlive">
            <div>{{ player.state.health }}</div>
            <div>
              <img src="/csgo-icons/weapon_armor_helmet.png" v-if="player.state.armor >  0 && player.state.helmet">
              <img src="/csgo-icons/weapon_armor.png" v-else-if="player.state.armor > 0">
            </div>
          </div>
          <img v-else src="/csgo-icons/skull-solid.svg" class="svg-icon"/>
        </div>
        <div class="right-box">
          <div class="health" :style="`background: linear-gradient(to left, black ${100-player.state.health}%, ${player.team.side === 'T' ? '#f52a2a' : '#3b5fed'} 0%);`">
            <div>
              {{ player.name }}
              <img
                v-if="isAlive && this.mainWeapon !== null"
                :src="`/csgo-icons/${(this.mainWeapon.state === 'active' ? this.mainWeapon.name : this.mainWeapon.name+'_holstered')}.png`"
                class="weapon-image"
              />
            </div>
          </div>
          <div>
            <div>
              <img src="/csgo-icons/crosshairs-solid.svg" class="svg-icon"/>
              {{player.stats.kills}}
              <img src="/csgo-icons/skull-solid.svg" class="svg-icon"/>
              {{player.stats.deaths}}
            </div>
            <div>
              $ {{player.state.money}}
            </div>
            <span v-if="isAlive">
              <span v-for="(weapon, key) in player.weapons" :key="key">
                <img
                  v-if="weapon.name !== mainWeapon.name && weapon.type !== 'Knife'"
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
.bordered {
  border: 5px solid rgb(80, 221, 80)
}
.hidden-border {
  border: 5px solid rgba(0,0,0,0)
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
  max-width: 70px;
}
.svg-icon {
  filter: invert(100%) sepia(5%) saturate(2878%) hue-rotate(0deg) brightness(100%) contrast(100%);
  height: 20px;
  width: 20px;
}
</style>

<script>
const MAIN_WEAPONS = ['Rifle', 'Submachine Gun', 'SniperRifle', 'Shotgun']

export default {
  props: ['player', 'focused'],
  data() {
    return {}
  },
  mounted() {},
  computed: {
    isAlive() {
      return this.player.state.health > 0
    },
    mainWeapon() {
      const weapons = {main: null, pistol: null}

      Object.entries(this.player.weapons).forEach(([, weapon]) => {
        if (MAIN_WEAPONS.includes(weapon.type))
          weapons.main = weapon
        else if (weapon.type === 'Pistol')
          weapons.pistol = weapon
      })

      if (weapons.main) return weapons.main
      if (weapons.pistol) return weapons.pistol
      
      return null
    },
  }
}
</script>