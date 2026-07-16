export async function onRequestGet(ctx) {
    if (ctx.params.sig.includes(".json")) {
        return await env.ASSETS.fetch(ctx.request);
    }
    return await env.ASSETS.fetch("/siegfried/update/" + ctx.params.sig + ".json");
}