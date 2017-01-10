/* eslint-disable no-unused-vars */

import Menu from './components/Menu.html'
import Home from './components/Home.html'
import Stacks from './components/Stacks.html'
import Topics from './components/Topics.html'
import StackEdit from './components/StackEdit.html'
import Functions from './components/Functions.html'
import KV from './components/Storage.html'

const sections = { Home, Stacks, Topics, StackEdit, Functions, KV }

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
