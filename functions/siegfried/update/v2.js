export async function onRequestGet(ctx) {
    return env.ASSETS.fetch("/siegfried/update/v2/update.json");
}