export async function onRequestGet(ctx) {
    return await env.ASSETS.fetch("/siegfried/update/v2/update.json");
}