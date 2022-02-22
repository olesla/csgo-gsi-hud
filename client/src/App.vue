<template>
  <div class="flex">
    <ul>
      <li v-for="player of terrorists" :key="player">
        <Player :player="player" :focused="state.player.steamid === player.steamid"/>
      </li>
    </ul>

    <div v-if="state && state.player">
      <focused-player :player="state.player"/>
    </div>
    
    <ul>
      <li v-for="player of counterTerrorists" :key="player">
        <Player :player="player" :focused="state.player.steamid === player.steamid"/>
      </li>
    </ul>
  </div>
</template>

<style scoped>
  ul {
    list-style-type: none;
  }
  .flex {
    display: flex;
    justify-content: space-between;
  }
</style>

<script>
import Player from './components/Player.vue'
import FocusedPlayer from './components/FocusedPlayer.vue'
import {CSGOGSI} from 'csgogsi'

const hasOwnProperty = (object, property) => Object.prototype.hasOwnProperty.call(object, property)

export default {
  name: 'App',
  components: {
    Player, FocusedPlayer
  },

  data() {
    return {
      socket: null,
      state: null,
    }
  },

  mounted() {
    const socket = new WebSocket("ws://localhost:9090/socket")
    const GSI = new CSGOGSI()
    this.socket = socket

    GSI.on('roundEnd', team => {
      console.log(`Team  ${team.name} win!`);
    });

    socket.addEventListener('open', () => {
      console.log('websocket connected')
    })

    socket.addEventListener("message", event => {
      const data = JSON.parse(event.data)
      const state = GSI.digest(data)

      if (this.state === null) {
        return this.state = state
      }      
      this.mergeState(state)
    })
    
  },

  methods: {
    team(side) {
      if (!this.state) return []
      const players = []
      for (const player of this.state.players) {
        if (player.team.side === side) {
          players.push(player)
        }
      }
      return players
    },
    mergeState(state) {

      for (const key of Object.keys(state)) {
        if (key === 'players') {
          continue
        }
        this.state[key] = state[key]
      }
      
      for (const [i, player] of this.state.players.entries()) {
        player.state = state.players[i].state
        // dont care about dead people
        if (player.state.health <= 0) { continue }

        player.stats = state.players[i].stats
        player.team = state.players[i].team

        // update weapons
        for (let j = 0; j < 8; j++) {
          const weaponId = `weapon_${j}`

          // weapon was used or thrown
          if (hasOwnProperty(this.state.players[i].weapons, weaponId) && !hasOwnProperty(state.players[i].weapons, weaponId)) {
            delete this.state.players[i].weapons[weaponId]
            continue
          }

          // weapon was added
          if (!hasOwnProperty(this.state.players[i].weapons, weaponId) && hasOwnProperty(state.players[i].weapons, weaponId)) {
            this.state.players[i].weapons[weaponId] = state.players[i].weapons[weaponId]
            continue
          }

          // update weapon
          if (hasOwnProperty(this.state.players[i].weapons, weaponId) && hasOwnProperty(state.players[i].weapons, weaponId)) {
            this.state.players[i].weapons[weaponId].name = state.players[i].weapons[weaponId].name
            this.state.players[i].weapons[weaponId].paintkit = state.players[i].weapons[weaponId].paintkit
            this.state.players[i].weapons[weaponId].state = state.players[i].weapons[weaponId].state
            this.state.players[i].weapons[weaponId].type = state.players[i].weapons[weaponId].type
            continue
          }
        }
      }
    
    }
  },

  computed: {
    counterTerrorists() {
      return this.team('CT')
    },
    terrorists() {
      return this.team('T')
    }
  }
}
</script>