import { map } from 'nanostores';
import type { App, BuildMetadata, NanoBuild } from '../../types/types';
import { useFetchBuild } from './api';

export type DetailsPageStoreType = {
	appId?: number;
	currentApp?: App;

	selectedBuildId?: number;

	currentBuild?: NanoBuild | undefined;

	availableBuilds?: BuildMetadata[];
};

export const detailsPageStore = map<DetailsPageStoreType>({});

detailsPageStore.listen((val, key) => {
	if (key === 'selectedBuildId') {
		useFetchBuild(val.selectedBuildId);
	}
});
