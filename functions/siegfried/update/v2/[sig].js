export async function onRequestGet(ctx) {
    if (ctx.params.sig.endsWith(".json")) {
        return await ctx.env.ASSETS.fetch(ctx.request);
    }
    const url = new URL("/siegfried/update/v2/" + ctx.params.sig + ".json", ctx.request.url);
    return await ctx.env.ASSETS.fetch(url);
}