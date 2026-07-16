export async function onRequestGet(ctx) {
    if (ctx.params.sig.includes(".json")) {
        return env.ASSETS.fetch(ctx.request.url);
    }
    return env.ASSETS.fetch("/siegfried/update/" + ctx.params.sig + ".json");
}