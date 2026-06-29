import { createI18n } from 'vue-i18n'

const messages = {
  en: {
    app: {
      title: 'NasFlow',
      tagline: 'Personal Task Management'
    },
    nav: {
      dashboard: 'Dashboard',
      tasks: 'Tasks',
      settings: 'Settings',
      logout: 'Logout'
    }
  },
  th: {
    app: {
      title: 'NasFlow',
      tagline: 'จัดการงานส่วนตัว'
    },
    nav: {
      dashboard: 'แดชบอร์ด',
      tasks: 'งาน',
      settings: 'การตั้งค่า',
      logout: 'ออกจากระบบ'
    }
  }
}

const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem('locale') || 'en',
  globalInjection: true,
  messages
})

export default i18n
