/* eslint-disable no-unused-vars */

import Menu from './components/Menu.html'
import Home from './components/Home.html'
import Stacks from './components/Stacks.html'
import Topics from './components/Topics.html'
import StackEdit from './components/StackEdit.html'

const sections = { Home, Stacks, Topics, StackEdit }

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
