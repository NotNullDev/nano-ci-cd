import { z } from 'zod';

export const AppSchema = z.object({
	ID: z.number(),
	CreatedAt: z.string(),
	UpdatedAt: z.string(),
	DeletedAt: z.string().nullable(),
	appName: z.string(),
	appStatus: z.string(),
	envVal: z.string(),
	envMountPath: z.string(),
	buildVal: z.string(),
	buildValMountPath: z.string(),
	repoUrl: z.string(),
	repoBranch: z.string()
});
export type App = z.infer<typeof AppSchema>;

export const NanoConfigSchema = z.object({
	globalEnvironment: z.string(),
	token: z.string()
});
export type NanoConfig = z.infer<typeof NanoConfigSchema>;

export const AppInfoSchema = z.object({
	appId: z.number(),
	appName: z.string(),
	appStatus: z.string(),
	envVal: z.string(),
	envMountPath: z.string(),
	buildVal: z.string(),
	buildValMountPath: z.string(),
	repoUrl: z.string()
});
export type AppInfo = z.infer<typeof AppInfoSchema>;

export const AppLogsTypeSchema = z.object({
	logs: z.string(),
	appId: z.number(),
	ID: z.number(),
	startedAt: z.string(),
	finishedAt: z.string(),
	buildStatus: z.string()
});
export type AppLogsType = z.infer<typeof AppLogsTypeSchema>;

export const NanoContextSchema = z.object({
	apps: z.array(AppSchema),
	nanoConfig: NanoConfigSchema,
	buildingAppId: z.number()
});

export type NanoContext = z.infer<typeof NanoContextSchema>;

export const NanoBuildSchema = z.object({
	ID: z.number(),
	appId: z.number(),
	buildStatus: z.string(),
	startedAt: z.string(),
	finishedAt: z.string(),
	logs: z.string()
});

export type NanoBuild = z.infer<typeof NanoBuildSchema>;

export const BuildMetadata = z.object({
	id: z.number(),
	date: z.string()
});

export type BuildMetadata = z.infer<typeof BuildMetadata>;
