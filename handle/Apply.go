package handle

import (
	"fmt"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/infrasonar/infrasonar-cli/cli"
	"github.com/infrasonar/infrasonar-cli/handle/util"
	"github.com/infrasonar/infrasonar-cli/req"
)

type Change struct {
	info string
	task any
}

type TaskUpsertZone struct {
	zone *cli.Zone
}

type TaskSetCollectorDisplay struct {
	collectorKey string
	display      bool
}

type TaskCreateAsset struct {
	asset *cli.AssetCli
}

type TaskSetAssetName struct {
	asset *cli.AssetCli
}

type TaskSetAssetMode struct {
	asset *cli.AssetCli
}

type TaskSetAssetKind struct {
	asset *cli.AssetCli
}

type TaskSetAssetZone struct {
	asset *cli.AssetCli
}

type TaskSetAssetDescription struct {
	asset *cli.AssetCli
}

type TaskAddLabelToAsset struct {
	asset *cli.AssetCli
	label *cli.Label
}

type TaskDeleteLabelFromAsset struct {
	asset *cli.AssetCli
	label *cli.Label
}

type TaskEnableAssetCheck struct {
	asset        *cli.AssetCli
	collectorKey string
	checkKey     string
}

type TaskDisableAssetCheck struct {
	asset        *cli.AssetCli
	collectorKey string
	checkKey     string
}

type TaskUpsertCollectorToAsset struct {
	asset        *cli.AssetCli
	collectorKey string
	config       map[string]any
}

type TaskRemoveCollectorFromAsset struct {
	asset        *cli.AssetCli
	collectorKey string
}

type TaskCreateLabel struct {
	label *cli.Label
}

type TaskSetLabelName struct {
	label *cli.Label
}

type TaskSetLabelColor struct {
	label *cli.Label
}

type TaskSetLabelDescription struct {
	label *cli.Label
}

func processChanges(api, token string, containerId int, changes *[]*Change) {
	n := len(*changes)
	for i, c := range *changes {
		fmt.Printf("Processing task %d/%d: %s ...\n", i+1, n, c.info)
		var err error
		switch task := c.task.(type) {
		case TaskUpsertZone:
			err = req.UpsertZone(api, token, containerId, task.zone.Zone, task.zone.Name)
		case TaskSetCollectorDisplay:
			err = req.SetCollectorDisplay(api, token, containerId, task.collectorKey, task.display)
		case TaskCreateAsset:
			task.asset.Id, err = req.CreateAsset(api, token, containerId, task.asset.Name)
		case TaskSetAssetName:
			err = req.SetAssetName(api, token, task.asset.Id, task.asset.Name)
		case TaskSetAssetMode:
			err = req.SetAssetMode(api, token, task.asset.Id, task.asset.Mode, nil)
		case TaskSetAssetKind:
			err = req.SetAssetKind(api, token, task.asset.Id, task.asset.Kind)
		case TaskSetAssetZone:
			err = req.SetAssetZone(api, token, task.asset.Id, *task.asset.Zone)
		case TaskSetAssetDescription:
			err = req.SetAssetDescription(api, token, task.asset.Id, task.asset.Description)
		case TaskAddLabelToAsset:
			err = req.AddLabelToAsset(api, token, task.asset.Id, task.label.Id)
		case TaskDeleteLabelFromAsset:
			err = req.DeleteLabelFromAsset(api, token, task.asset.Id, task.label.Id)
		case TaskEnableAssetCheck:
			err = req.EnableAssetCheck(api, token, task.asset.Id, task.collectorKey, task.checkKey)
		case TaskDisableAssetCheck:
			err = req.DisableAssetCheck(api, token, task.asset.Id, task.collectorKey, task.checkKey)
		case TaskUpsertCollectorToAsset:
			err = req.UpsertCollectorToAsset(api, token, task.asset.Id, task.collectorKey, task.config)
		case TaskRemoveCollectorFromAsset:
			err = req.RemoveCollectorFromAsset(api, token, task.asset.Id, task.collectorKey)
		case TaskCreateLabel:
			task.label.Id, err = req.CreateLabel(api, token, containerId, task.label.Name)
		case TaskSetLabelName:
			err = req.SetLabelName(api, token, task.label.Id, task.label.Name)
		case TaskSetLabelColor:
			err = req.SetLabelColor(api, token, task.label.Id, task.label.Color)
		case TaskSetLabelDescription:
			err = req.SetLabelDescription(api, token, task.label.Id, task.label.Description)

		}
		util.ExitOnErr(err)
	}
}

func cval(a any) string {
	return color.HiCyanString("%v", a)
}

func readLabelChanges(changes *[]*Change, cl, tl *cli.Label) {
	if cl.Id != tl.Id {
		panic("label ID mismatch")
	}
	if tl.Name != "" && cl.Name != "" && tl.Name != cl.Name {
		*changes = append(*changes, &Change{
			info: fmt.Sprintf("Set name for label ID %s to: %s", cval(tl.Id), cval(tl.Name)),
			task: TaskSetLabelName{label: tl},
		})
	}
	if tl.Color != "" && tl.Color != cl.Color {
		*changes = append(*changes, &Change{
			info: fmt.Sprintf("Set color for label ID %s to '%s'", cval(tl.Id), cval(tl.Color)),
			task: TaskSetLabelColor{label: tl},
		})
	}
	if tl.Description != "" && tl.Description != cl.Description {
		*changes = append(*changes, &Change{
			info: fmt.Sprintf("Set description for label ID %s to '%s'", cval(tl.Id), cval(util.Short(tl.Description, 12))),
			task: TaskSetLabelDescription{label: tl},
		})
	}
}

func assetChanges(changes *[]*Change, purge bool, ca, ta *cli.AssetCli, cs, ts *cli.State) {
	if ta.Name != "" && ca.Name != "" && ta.Name != ca.Name {
		*changes = append(*changes, &Change{
			info: fmt.Sprintf("Set name for asset '%s' to: '%s'", cval(ta.Str()), cval(ta.Name)),
			task: TaskSetAssetName{asset: ta},
		})
	}
	if ta.Mode != "" && ta.Mode != ca.Mode {
		*changes = append(*changes, &Change{
			info: fmt.Sprintf("Set mode for asset '%s' to: '%s'", cval(ta.Str()), cval(ta.Mode)),
			task: TaskSetAssetMode{asset: ta},
		})
	}
	if ta.Kind != "" && ta.Kind != ca.Kind {
		*changes = append(*changes, &Change{
			info: fmt.Sprintf("Set kind for asset '%s' to: '%s'", cval(ta.Str()), cval(ta.Kind)),
			task: TaskSetAssetKind{asset: ta},
		})
	}
	if ta.Zone != nil && (ca.Zone == nil || *ta.Zone != *ca.Zone) {
		*changes = append(*changes, &Change{
			info: fmt.Sprintf("Set zone for asset '%s' to: %s", cval(ta.Str()), cval(*ta.Zone)),
			task: TaskSetAssetZone{asset: ta},
		})
	}
	if ta.Description != "" && ta.Description != ca.Description {
		*changes = append(*changes, &Change{
			info: fmt.Sprintf("Set description for asset '%s' to: '%s'", cval(ta.Str()), cval(util.Short(ta.Description, 12))),
			task: TaskSetAssetDescription{asset: ta},
		})
	}
	if ta.Labels != nil {
		for _, key := range *ta.Labels {
			if label := ts.LabelByKey(key); label != nil {
				if !ca.HasLabelId(label.Id, cs.GetLabelMap()) {
					*changes = append(*changes, &Change{
						info: fmt.Sprintf("Add label '%s' to asset '%s'", cval(label.Str()), cval(ta.Str())),
						task: TaskAddLabelToAsset{asset: ta, label: label},
					})
				}
			}
		}
	}
	if ta.Collectors != nil {
		for _, collector := range *ta.Collectors {
			var other *cli.TCollector

			if ca.Collectors != nil {
				for _, c := range *ca.Collectors {
					if c.Key == collector.Key {
						other = &c
						break
					}
				}
			}

			if other == nil {
				*changes = append(*changes, &Change{
					info: fmt.Sprintf("Add collector '%s' to asset '%s'", cval(collector.Key), cval(ta.Str())),
					task: TaskUpsertCollectorToAsset{asset: ta, collectorKey: collector.Key, config: collector.Config},
				})
			} else {
				for k, v := range collector.Config {
					ov, ok := other.Config[k]
					if !ok || !reflect.DeepEqual(v, ov) {
						*changes = append(*changes, &Change{
							info: fmt.Sprintf("Update collector '%s' configuration for asset '%s'", cval(collector.Key), cval(ta.Str())),
							task: TaskUpsertCollectorToAsset{asset: ta, collectorKey: collector.Key, config: collector.Config},
						})
						break
					}
				}
			}
		}
	}
	if ta.DisabledChecks != nil {
		for _, disabledChk := range *ta.DisabledChecks {
			found := false
			if ca.DisabledChecks != nil {
				for _, c := range *ca.DisabledChecks {
					if c.Collector == disabledChk.Collector && c.Check == disabledChk.Check {
						found = true
						break
					}
				}
			}
			if !found {
				*changes = append(*changes, &Change{
					info: fmt.Sprintf("Disable collector check '%s/%s' on asset '%s'", cval(disabledChk.Collector), cval(disabledChk.Check), cval(ta.Str())),
					task: TaskDisableAssetCheck{asset: ta, collectorKey: disabledChk.Collector, checkKey: disabledChk.Check},
				})
			}
		}
	}

	if purge {
		if ca.Labels != nil && ta.Labels != nil {
			for _, key := range *ca.Labels {
				if label := cs.LabelByKey(key); label != nil {
					if !ta.HasLabelId(label.Id, ts.GetLabelMap()) {
						*changes = append(*changes, &Change{
							info: fmt.Sprintf("Delete label '%s' from asset '%s'", cval(label.Str()), cval(ta.Str())),
							task: TaskDeleteLabelFromAsset{asset: ta, label: label},
						})
					}
				}
			}
		}
		if ca.DisabledChecks != nil && ta.DisabledChecks != nil {
			for _, disabledChk := range *ca.DisabledChecks {
				found := false
				for _, c := range *ta.DisabledChecks {
					if c.Collector == disabledChk.Collector && c.Check == disabledChk.Check {
						found = true
						break
					}
				}
				if !found {
					*changes = append(*changes, &Change{
						info: fmt.Sprintf("Enable collector check '%s/%s' on asset '%s'", cval(disabledChk.Collector), cval(disabledChk.Check), cval(ta.Str())),
						task: TaskEnableAssetCheck{asset: ta, collectorKey: disabledChk.Collector, checkKey: disabledChk.Check},
					})
				}
			}
		}
		if ca.Collectors != nil && ta.Collectors != nil {
			for _, collector := range *ca.Collectors {
				found := false
				for _, c := range *ta.Collectors {
					if c.Key == collector.Key {
						found = true
						break
					}
				}
				if !found {
					*changes = append(*changes, &Change{
						info: fmt.Sprintf("Remove collector '%s' from asset '%s'", cval(collector.Key), cval(ta.Str())),
						task: TaskRemoveCollectorFromAsset{asset: ta, collectorKey: collector.Key},
					})
				}
			}
		}
	}
}

func ensureCollectors(ts *cli.State, cMap map[string]*cli.Collector) []*Change {
	changes := []*Change{}

	//
	// Show collectors and sanity check mode and disabled checks
	//
	enableCollector := cli.StrSet{}
	for _, ta := range ts.Assets {
		if ta.Collectors != nil {
			for _, c := range *ta.Collectors {
				if enableCollector.Has(c.Key) {
					continue
				}
				if _, ok := cMap[c.Key]; !ok {
					changes = append(changes, &Change{
						info: fmt.Sprintf("Enable collector: %s", cval(c.Key)),
						task: TaskSetCollectorDisplay{collectorKey: c.Key, display: true},
					})
					enableCollector.Set(c.Key)
				}
			}
		}
	}
	return changes
}

func ensureChanges(api, token string, purge bool, cs, ts *cli.State, cMap map[string]*cli.Collector) []*Change {
	changes := []*Change{}
	//
	// Container changes
	//
	// if ts.Container.Name != "" && ts.Container.Name != cs.Container.Name {
	// 	changes = append(changes, &Change{
	// 		info: fmt.Sprintf("Rename container ID %s from '%s' to '%s'", cval(cs.Container.Id), cval(cs.Container.Name), cval(ts.Container.Name)),
	// 		// TODO: Set name container API
	// 	})
	// }
	//
	// Show collectors and sanity check mode and disabled checks
	//
	enableCollector := cli.StrSet{}
	for _, ta := range ts.Assets {
		if ta.Collectors != nil {
			for _, c := range *ta.Collectors {
				if enableCollector.Has(c.Key) {
					continue
				}
				if _, ok := cMap[c.Key]; !ok {
					changes = append(changes, &Change{
						info: fmt.Sprintf("Enable collector: %s", cval(c.Key)),
						task: TaskSetCollectorDisplay{collectorKey: c.Key, display: true},
					})
					enableCollector.Set(c.Key)
				}
			}
		}
		if ta.DisabledChecks != nil {
			for _, disabledChk := range *ta.DisabledChecks {
				found := false
				if ta.Collectors != nil {
					for _, c := range *ta.Collectors {
						if c.Key == disabledChk.Collector {
							found = true
							break
						}
					}
					if !found {
						util.ExitErr("Collector '%s' is not configured for asset '%s', but a disabled check for it exists.", disabledChk.Collector, ta.Str())
					}
				}
			}
		}
		switch ta.Mode {
		case "", "normal", "maintenance", "disabled":
			continue
		}
		util.ExitErr("Asset '%s' has an invalid mode '%s'. Must be one of {normal,maintenance,disabled}", ta.Str(), ta.Mode)
	}

	//
	// Zone changes
	//
	for _, tz := range ts.Zones {
		cz := cs.ZoneById(tz.Zone)
		if cz == nil {
			if tz.Zone < 1 || tz.Zone > 9 {
				util.ExitErr("Invalid zone '%d'. Must be a value between 1 and 9.", tz.Zone)
			}
			if tz.Name == "" {
				util.ExitErr("Zone '%d' is new and therefore requires a name", tz.Zone)
			}
			changes = append(changes, &Change{
				info: fmt.Sprintf("Create new zone: %s", cval(tz.Str())),
				task: TaskUpsertZone{zone: tz},
			})
		} else {
			if tz.Name != "" && tz.Name != cz.Name {
				changes = append(changes, &Change{
					info: fmt.Sprintf("Rename zone ID %s from '%s' to '%s'", cval(tz.Zone), cval(cz.Name), cval(tz.Name)),
					task: TaskUpsertZone{zone: tz},
				})
			}
		}
	}

	missingLabelIds := cli.IntSet{}
	for _, tl := range ts.Labels {
		if tl.Id == 0 {
			// New label
			if tl.Name == "" {
				util.ExitErr("One or more labels are missing both an 'id' and a 'name'. At least one of these attributes is required for each label.")
			}
			changes = append(changes, &Change{
				info: fmt.Sprintf("Create new label: %s", cval(tl.Name)),
				task: TaskCreateLabel{label: tl},
			})
			readLabelChanges(&changes, &cli.DefaultLabel, tl)
		} else {
			// Label with ID
			if cl := cs.LabelById(tl.Id); cl == nil {
				missingLabelIds.Set(tl.Id)
			} else {
				readLabelChanges(&changes, cl, tl)
			}
		}
	}
	if len(missingLabelIds) > 0 {
		fmt.Println("Get missing labels...")
		lm, err := req.GetLabels(api, token, missingLabelIds)
		util.ExitOnErr(err)
		for _, tl := range ts.Labels {
			if tl.Id != 0 {
				if cl := cs.LabelById(tl.Id); cl == nil {
					if cl := lm.LabelById(tl.Id); cl != nil {
						readLabelChanges(&changes, cl, tl)
					}
				}
			}
		}
	}
	for _, ta := range ts.Assets {
		if ta.Labels != nil {
			for _, labelKey := range *ta.Labels {
				if _, ok := ts.Labels[labelKey]; !ok {
					util.ExitErr("Asset '%s' is using label reference '%s' which does not exist in 'labels'.", ta.Str(), labelKey)
				}
			}
		}
		if ta.Zone != nil {
			if zone := ts.ZoneById(*ta.Zone); zone == nil {
				util.ExitErr("Asset '%s' is using zone ID %d which does not exist in 'zones'.", ta.Str(), *ta.Zone)
			}
		}

		if ta.Id == 0 {
			// New asset
			if ta.Name == "" {
				util.ExitErr("One or more assets are missing both an 'id' and a 'name'. At least one of these attributes is required for each asset.")
			}
			changes = append(changes, &Change{
				info: fmt.Sprintf("Create new asset: %s", cval(ta.Name)),
				task: TaskCreateAsset{asset: ta},
			})
			assetChanges(&changes, purge, &cli.DefaultAsset, ta, cs, ts)
		} else {
			ca := cs.AssetById(ta.Id)
			if ca == nil {
				util.ExitErr("Asset ID %d not found in container '%s'.", ta.Id, cs.Container.Str())
			}
			assetChanges(&changes, purge, ca, ta, cs, ts)
		}
	}
	return changes
}

func getCacheState(containerId int) *cli.State {
	state := cli.StateFromCache(containerId)
	if state != nil {
		if age, err := state.GetAge(); err == nil {
			if age.Hours() < 8 {
				util.Color("A cache for container ID %d was found that is only %s old. Would you like to use it? (yes/no): ", containerId, util.HumanizeDuration(*age))
				if util.AskForConfirmation() {
					return state
				}
			}
		}
	}
	return nil
}

func sanitizeConfig(input map[string]any) map[string]any {
	clone := make(map[string]any)
	for k, v := range input {
		if k == "password" || k == "secret" {
			clone[k] = "xxx"
			continue
		}
		clone[k] = v
	}
	return clone
}

func sanityCheckCollectorConfig(ts *cli.State, cMap map[string]*cli.Collector, api, token string, remoteValidation bool) {
	for _, asset := range ts.Assets {
		if asset.Collectors == nil {
			continue
		}
		for _, collector := range *asset.Collectors {
			if c, ok := cMap[collector.Key]; ok {
				for k, v := range collector.Config {
					found := false
					for _, o := range c.Options {
						if o.Key == k {
							found = true
							if o.Type == "String" && (k == "password" || k == "secret") {
								if _, ok := v.(string); ok {
									break
								}
								if obj, ok := v.(map[string]any); ok {
									if v, ok := obj["encrypted"]; ok && len(obj) == 1 {
										if _, ok := v.(string); ok {
											break
										}
									}
								}
								util.ExitErr("Collector '%s' on asset '%s' expects property '%s' to be a string or encryption value", collector.Key, asset.Str(), k)
							}
							switch o.Type {
							case "Bool":
								if _, ok := v.(bool); !ok {
									util.ExitErr("Collector '%s' on asset '%s' expects a boolean value for property '%s' but found type %T", collector.Key, asset.Str(), k, v)
								}
							case "Int":
								if _, ok := v.(int); !ok {
									util.ExitErr("Collector '%s' on asset '%s' expects an integer value for property '%s' but found type %T", collector.Key, asset.Str(), k, v)
								}
							case "Float":
								if _, ok := v.(float64); !ok {
									util.ExitErr("Collector '%s' on asset '%s' expects a floating point for property '%s' but found type %T", collector.Key, asset.Str(), k, v)
								}
							case "String":
								if _, ok := v.(string); !ok {
									util.ExitErr("Collector '%s' on asset '%s' expects a string value for property '%s' but found type %T", collector.Key, asset.Str(), k, v)
								}
							case "ListBool", "ListInt", "ListFloat", "ListString":
								if arr, ok := v.([]any); ok {
									switch o.Type {
									case "ListBool":
										for _, v := range arr {
											if _, ok = v.(bool); !ok {
												util.ExitErr("Collector '%s' on asset '%s' expects a list of boolean values for property '%s' but the list contains type %T", collector.Key, asset.Str(), k, v)
											}
										}
									case "ListInt":
										for _, v := range arr {
											if _, ok = v.(int); !ok {
												util.ExitErr("Collector '%s' on asset '%s' expects a list of integer values for property '%s' but the list contains type %T", collector.Key, asset.Str(), k, v)
											}
										}
									case "ListFloat":
										for _, v := range arr {
											if _, ok = v.(float64); !ok {
												util.ExitErr("Collector '%s' on asset '%s' expects a list of floating point values for property '%s' but the list contains type %T", collector.Key, asset.Str(), k, v)
											}
										}
									case "ListString":
										for _, v := range arr {
											if _, ok = v.(string); !ok {
												util.ExitErr("Collector '%s' on asset '%s' expects a list of string values for property '%s' but the list contains type %T", collector.Key, asset.Str(), k, v)
											}
										}
									}
									break
								}
								util.ExitErr("Collector '%s' on asset '%s' expects a list of values for property '%s' but found type %T", collector.Key, asset.Str(), k, v)
							}
							break
						}
					}
					if !found {
						util.ExitErr("Collector '%s' on asset '%s' contains an unknown configuration property '%s'.", collector.Key, asset.Str(), k)
					}
				}
			}
			if remoteValidation {
				err := req.VerifyCollectorConfig(api, token, collector.Key, sanitizeConfig(collector.Config))
				if err != nil {
					util.ExitErr("Collector '%s' on asset '%s': %s", collector.Key, asset.Str(), err)
				}
			}
		}
	}
}

func revertUse(assets []*cli.AssetCli) {
	for _, asset := range assets {
		if asset.Collectors != nil {
			for _, c := range *asset.Collectors {
				if v, ok := c.Config["use"]; ok {
					delete(c.Config, "use")
					c.Config["_use"] = v
				}
			}
		}
	}
}

func ensureNumbersAndDefaults(assets []*cli.AssetCli, cMap map[string]*cli.Collector) {
	for _, asset := range assets {
		if asset.Collectors == nil {
			continue
		}
		for _, collector := range *asset.Collectors {
			if c, ok := cMap[collector.Key]; ok {
				for k, v := range collector.Config {
					for _, o := range c.Options {
						if o.Key == k {
							switch o.Type {
							case "Int":
								if v, ok := v.(float64); ok {
									if util.IsIntegral(v) {
										collector.Config[k] = int(v)
									}
								}
							case "Float":
								if v, ok := v.(int); ok {
									collector.Config[k] = float64(v)
								}
							case "ListInt", "ListFloat":
								if arr, ok := v.([]any); ok {
									switch o.Type {
									case "ListInt":
										for i, v := range arr {
											if v, ok := v.(float64); ok {
												if util.IsIntegral(v) {
													arr[i] = int(v)
												}
											}
										}
									case "ListFloat":
										for i, v := range arr {
											if v, ok := v.(int); ok {
												arr[i] = float64(v)
											}
										}
									}
								}
							}
							break
						}
					}
				}
				for _, o := range c.Options {
					switch o.Type {
					case "Int":
						if v, ok := o.Default.(float64); ok {
							if util.IsIntegral(v) {
								o.Default = int(v)
							}
						}
					case "Float":
						if v, ok := o.Default.(int); ok {
							o.Default = float64(v)
						}
					case "ListInt", "ListFloat":
						if arr, ok := o.Default.([]any); ok {
							switch o.Type {
							case "ListInt":
								for i, v := range arr {
									if v, ok := v.(float64); ok {
										if util.IsIntegral(v) {
											arr[i] = int(v)
										}
									}
								}
							case "ListFloat":
								for i, v := range arr {
									if v, ok := v.(int); ok {
										arr[i] = float64(v)
									}
								}
							}
						}
					}
					if collector.Config != nil {
						if _, ok := collector.Config[o.Key]; !ok {
							collector.Config[o.Key] = o.Default
						}
					}
				}
			}
		}
	}
}

func niceAssetKinds(api string, assets []*cli.AssetCli) {
	kinds, err := req.GetAssetKinds(api)
	util.ExitOnErr(err)

	for _, asset := range assets {
		if asset.Kind != "" {
			kind := util.InSlice(kinds, asset.Kind)
			if kind == nil {
				util.ExitErr("Asset '%s' has an invalid asset kind: %s", asset.Str(), asset.Kind)
			}
			asset.Kind = *kind
		}
	}
}

func updateCollectorMap(api, token string, containerId int, cMap *map[string]*cli.Collector) error {
	collectors, err := req.GetCollectors(api, token, containerId, []string{"key"}, true)
	if err != nil {
		return err
	}

	for _, c := range collectors {
		(*cMap)[c.Key] = c
	}
	return nil
}

func Apply(api, token, filename string, dryRun, purge bool) {
	if dryRun {
		util.Color(`-----------------------------------------
  Simulation :: no changes will be made
-----------------------------------------
`)
	}
	fmt.Println("Read input file...")
	ts, err := cli.StateFromFile(filename)
	util.ExitOnErr(err)

	if ts.Container.Id == 0 {
		util.ExitErr("missing container ID in input file")
	}

	if !dryRun {
		fmt.Println("Check token permissions...")
		me, err := req.GetMe(api, token, ts.Container.Id)
		util.ExitOnErr(err)
		util.ExitOnErr(me.CheckApplyPermissions())
	}

	cs := getCacheState(ts.Container.Id)
	if cs == nil {
		fmt.Println("Read current state...")
		cs = ensureState(&TGetAssets{
			Api:             api,
			Token:           token,
			Output:          "",
			OutFn:           "-", // Force progress output, nothing will be written
			Container:       ts.Container.Id,
			Asset:           0,
			Properties:      cli.AssetProperties,
			Filters:         []string{},
			IncludeDefaults: true,
		})
		cs.WriteCache()
	}

	cMap := map[string]*cli.Collector{}

	if ts.HasAssetKind() {
		niceAssetKinds(api, ts.Assets)
	}

	if ts.HasCollector() {
		// use -> _use
		revertUse(cs.Assets)
		revertUse(ts.Assets)

		fmt.Println("Read collectors...")
		util.ExitOnErr(updateCollectorMap(api, token, ts.Container.Id, &cMap))

		changes := ensureCollectors(ts, cMap)
		n := len(changes)
		if n > 0 {
			if dryRun {
				util.Color("To run a more accurate dry run, %d collector%s need to be enabled. Proceed? (yes/no): ", n, util.Plural(n))
			} else {
				util.Color("To continue, %d collector%s must be enabled. Proceed? (yes/no): ", n, util.Plural(n))
			}
			if util.AskForConfirmation() {
				ts.ClearCache() // Clear the cache as we're about to make changes
				fmt.Println("")
				processChanges(api, token, ts.Container.Id, &changes)
				fmt.Println("")
			} else {
				util.ExitOk("Cancelled.")
			}
			// We need this sleep, for the collector changes
			time.Sleep(1 * time.Second)

			fmt.Println("Read collectors...")
			util.ExitOnErr(updateCollectorMap(api, token, ts.Container.Id, &cMap))
		}

		ensureNumbersAndDefaults(cs.Assets, cMap)
		ensureNumbersAndDefaults(ts.Assets, cMap)

		util.Color("Perform remote configuration check? (Local check is faster, remote is more thorough) (yes/no): ")
		remoteValidation := util.AskForConfirmation()
		fmt.Println("")

		sanityCheckCollectorConfig(ts, cMap, api, token, remoteValidation)
	}

	changes := ensureChanges(api, token, purge, cs, ts, cMap)
	n := len(changes)

	if n == 0 {
		util.ExitOk("No changes found.")
	}

	util.Color("Found %d change%s. Show details? (yes/no): ", n, util.Plural(n))
	if util.AskForConfirmation() {
		fmt.Println("")
		for _, c := range changes {
			fmt.Printf("- %s\n", c.info)
		}
		fmt.Println("")
	}

	if dryRun {
		util.ExitOk("Done. (no changes made)")
	} else {
		util.Color("Do you want to apply the change%s? (yes/no): ", util.Plural(n))
		if util.AskForConfirmation() {
			ts.ClearCache() // Clear the cache as we're about to make changes
			fmt.Println("")
			processChanges(api, token, ts.Container.Id, &changes)
			fmt.Println("")
			util.ExitOk("Done.")
		}
		util.ExitOk("Cancelled.")
	}
}
