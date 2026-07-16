export async function onRequestGet(ctx) {
    if (ctx.params.sig.includes(".json")) {
        return await ctx.env.ASSETS.fetch(ctx.request);
    }
    return await ctx.env.ASSETS.fetch("/siegfried/update/v2/" + ctx.params.sig + ".json");
}