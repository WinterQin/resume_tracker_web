import api from './index'

export const applicationApi = {
  getApplications(params) {
    return api.get('/applications', { params })
  },
  createApplication(data) {
    return api.post('/applications', data)
  },
  updateApplication(id, data) {
    return api.put(`/applications/${id}`, {
      id: parseInt(id),
      ...data
    })
  },
  deleteApplication(id) {
    return api.delete(`/applications/${id}`)
  },
  updateStatus(data) {
    return api.patch('/applications/status', data)
  },
  updateEvent(data) {
    return api.put('/applications/event', data)
  },
  getStatistics() {
    return api.get('/applications/statistics')
  },
  getUpcomingEvents() {
    return api.get('/applications/upcoming-events')
  },
  getRecentApplications() {
    return api.get('/applications/recent')
  }
} 