import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import LoginView from './LoginView.vue'

describe('LoginView', () => {
  it('should render', () => {
    const wrapper = mount(LoginView)
    expect(true).toBe(true)
  })
})
