const getters = {
  sidebar: state => state.app.sidebar,
  device: state => state.app.device,
  token: state => state.user.token,
  avatar: state => state.user.avatar,
  name: state => state.user.name,
  userType: state => state.user.userType,
  dynamicAttributes: state => state.user.dynamicAttributes
}
export default getters
