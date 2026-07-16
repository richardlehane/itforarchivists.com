export async function onRequestGet(ctx) {
    if (ctx.params.sig.includes(".json")) {
        return env.ASSETS.fetch(ctx.request.url);
    }
    return env.ASSETS.fetch("/siegfried/update/v2/" + ctx.params.sig + ".json");
}