import type { PageLoad } from './$types';
import { env } from '$env/dynamic/public';

export const load: PageLoad = (async ({ params, fetch }) => {
  const slug = params.slug.split("-");
  const page_id = slug[slug.length-1]
  const res = await fetch(`${env.PUBLIC_API_URL}/story/${page_id}`);
  const story = await res.json();
  return { story };
})
