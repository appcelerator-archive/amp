/* eslint-disable no-unused-vars */

import Menu from './components/Menu.html'
import Home from './components/Home.html'
import Stacks from './components/Stacks.html'
import Topics from './components/Topics.html'

const sections = { Home, Stacks, Topics }

const MenuComponent = new Menu({
  target: document.querySelector('#menu'),
})

let ActiveComponent = {
  teardown () {}
}

MenuComponent.observe('active', active => {
  ActiveComponent.teardown()
  ActiveComponent = new sections[active]({
    target: document.querySelector('main')
  })
})

// const api = new AmpApi()

// const LogDisplayComponent = new LogDisplay({
//   target: document.querySelector('.log-display'),
// })

// api.logs({
//   service: 'amplifier'
// }).then(logs => {
//   LogDisplayComponent.set({logs})
// })
