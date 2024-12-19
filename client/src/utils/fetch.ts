import router from '@/router'

// This method is a simple wrapper around fetch that checks for 401 and 403
// and redirects to the login page if the user is not authenticated.
export async function vfetch(
  url: RequestInfo | URL,
  options?: RequestInit,
): Promise<Response> {
  return await fetch(url, options).then(response => {
    if (response.status === 401 || response.status === 403) {
      router.push({ name: 'login' })
      throw new Error('Unauthorized')
    }
    return response
  })
}
