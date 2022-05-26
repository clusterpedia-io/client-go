package memory

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	internal "github.com/clusterpedia-io/api/clusterpedia"
	fields2 "github.com/clusterpedia-io/api/clusterpedia/fields"
	"github.com/clusterpedia-io/client-go/constants"
	"github.com/clusterpedia-io/client-go/fake/storage"
	"github.com/clusterpedia-io/client-go/fake/utils"
	"gorm.io/datatypes"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type Resource struct {
	ID uint `gorm:"primaryKey"`

	Group    string `gorm:"size:63;not null;uniqueIndex:uni_group_version_resource_cluster_namespace_name;index:idx_group_version_resource_namespace_name;index:idx_group_version_resource_name"`
	Version  string `gorm:"size:15;not null;uniqueIndex:uni_group_version_resource_cluster_namespace_name;index:idx_group_version_resource_namespace_name;index:idx_group_version_resource_name"`
	Resource string `gorm:"size:63;not null;uniqueIndex:uni_group_version_resource_cluster_namespace_name;index:idx_group_version_resource_namespace_name;index:idx_group_version_resource_name"`
	Kind     string `gorm:"size:63;not null"`

	Cluster         string    `gorm:"size:253;not null;uniqueIndex:uni_group_version_resource_cluster_namespace_name,length:100;index:idx_cluster"`
	Namespace       string    `gorm:"size:253;not null;uniqueIndex:uni_group_version_resource_cluster_namespace_name,length:50;index:idx_group_version_resource_namespace_name"`
	Name            string    `gorm:"size:253;not null;uniqueIndex:uni_group_version_resource_cluster_namespace_name,length:100;index:idx_group_version_resource_namespace_name;index:idx_group_version_resource_name"`
	OwnerUID        types.UID `gorm:"column:owner_uid;size:36;not null;default:''"`
	UID             types.UID `gorm:"size:36;not null"`
	ResourceVersion string    `gorm:"size:30;not null"`

	Object datatypes.JSON `gorm:"not null"`

	CreatedAt time.Time `gorm:"not null"`
	SyncedAt  time.Time `gorm:"not null;autoUpdateTime"`
	DeletedAt sql.NullTime
}

type FakeStorageFactory struct {
	storage map[schema.GroupVersionKind]*FakeResourceStorage
}

func NewFakeStorageFactory() *FakeStorageFactory {
	return &FakeStorageFactory{
		storage: make(map[schema.GroupVersionKind]*FakeResourceStorage),
	}
}

func (fake *FakeStorageFactory) GetResourceVersions(ctx context.Context, cluster string) (map[schema.GroupVersionResource]map[string]interface{}, error) {
	return nil, nil
}

func (fake *FakeStorageFactory) CleanCluster(ctx context.Context, cluster string) error {
	return nil
}

func (fake *FakeStorageFactory) CleanClusterResource(ctx context.Context, cluster string, gvr schema.GroupVersionResource) error {
	return nil
}

func (fake *FakeStorageFactory) Create(ctx context.Context, cluster string, obj runtime.Object) error {
	gvk := schema.GroupVersionKind{
		Group:   obj.GetObjectKind().GroupVersionKind().Group,
		Kind:    obj.GetObjectKind().GroupVersionKind().Kind,
		Version: obj.GetObjectKind().GroupVersionKind().Version,
	}
	resourceStorage, exist := fake.storage[gvk]
	if !exist {
		return fmt.Errorf("fail init fake resource storage")
	}

	return resourceStorage.Create(ctx, cluster, obj)
}

func (fake *FakeStorageFactory) NewResourceStorage(config *storage.ResourceStorageConfig) (storage.ResourceStorage, error) {
	gvk := schema.GroupVersionKind{
		Group:   config.GroupResource.Group,
		Version: config.StorageVersion.Version,
		Kind:    config.Kind,
	}
	if _, exist := fake.storage[gvk]; !exist {
		resourceStorage := &FakeResourceStorage{
			DB:                   NewDataBase(),
			codec:                config.Codec,
			storageGroupResource: config.StorageGroupResource,
			storageVersion:       config.StorageVersion,
			memoryVersion:        config.MemoryVersion,
		}
		fake.storage[gvk] = resourceStorage
	}

	return fake.storage[gvk], nil
}

func (fake *FakeStorageFactory) GetResourceStorage(gvk schema.GroupVersionKind) (storage.ResourceStorage, error) {
	if _, exist := fake.storage[gvk]; !exist {
		return nil, fmt.Errorf("fail init fake resource storage")
	}
	return fake.storage[gvk], nil
}

func (fake *FakeStorageFactory) NewCollectionResourceStorage(cr *internal.CollectionResource) (storage.CollectionResourceStorage, error) {
	return nil, nil
}

func (fake *FakeStorageFactory) GetCollectionResources(ctx context.Context) ([]*internal.CollectionResource, error) {
	return nil, nil
}

type FakeResourceStorage struct {
	DB    *DataBase
	codec runtime.Codec

	storageGroupResource schema.GroupResource
	storageVersion       schema.GroupVersion
	memoryVersion        schema.GroupVersion
}

func NewFakeResourceStorage() *FakeResourceStorage {
	config := &storage.ResourceStorageConfig{
		Codec: NewCodec(),
	}
	return &FakeResourceStorage{
		DB:                   NewDataBase(),
		codec:                config.Codec,
		storageGroupResource: config.StorageGroupResource,
		storageVersion:       config.StorageVersion,
		memoryVersion:        config.MemoryVersion,
	}
}

func (f *FakeResourceStorage) GetStorageConfig() *storage.ResourceStorageConfig {
	return nil
}

func (f *FakeResourceStorage) Get(ctx context.Context, cluster, namespace, name string, obj runtime.Object) error {
	index, err := f.GetIndex(ctx, cluster, namespace, name)
	if err != nil {
		return err
	}
	buf, err := f.DB.Table[index].Object.MarshalJSON()
	if err != nil {
		return err
	}
	_, _, err = f.codec.Decode(buf, nil, obj)
	if err != nil {
		return err
	}
	return nil
}

func (f *FakeResourceStorage) List(ctx context.Context, listObj runtime.Object, opts *internal.ListOptions) error {
	var err error
	data := f.DB.ListSimpleSearch(opts)
	objects := make([][]byte, len(data))
	for k, v := range data {
		objects[k], err = v.Object.MarshalJSON()
		if err != nil {
			return err
		}
	}

	listPtr, err := meta.GetItemsPtr(listObj)
	if err != nil {
		return err
	}

	v, err := conversion.EnforcePtr(listPtr)
	if err != nil || v.Kind() != reflect.Slice {
		return err
	}

	newItemFunc := getNewItemFunc(listObj, v)
	for _, object := range objects {
		if err := appendListItem(v, object, f.codec, newItemFunc); err != nil {
			return err
		}
	}
	return nil
}

func getNewItemFunc(listObj runtime.Object, v reflect.Value) func() runtime.Object {
	if _, isUnstructuredList := listObj.(*unstructured.UnstructuredList); isUnstructuredList {
		return func() runtime.Object {
			return &unstructured.Unstructured{Object: map[string]interface{}{}}
		}
	}

	elem := v.Type().Elem()
	return func() runtime.Object {
		return reflect.New(elem).Interface().(runtime.Object)
	}
}

func appendListItem(v reflect.Value, data []byte, codec runtime.Codec, newItemFunc func() runtime.Object) error {
	obj, _, err := codec.Decode(data, nil, newItemFunc())
	if err != nil {
		return err
	}
	v.Set(reflect.Append(v, reflect.ValueOf(obj).Elem()))
	return nil
}

func (f *FakeResourceStorage) Create(ctx context.Context, cluster string, obj runtime.Object) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	var ownerUID types.UID
	if owner := metav1.GetControllerOfNoCopy(accessor); owner != nil {
		ownerUID = owner.UID
	}

	var buf bytes.Buffer
	if err := f.codec.Encode(obj, &buf); err != nil {
		return err
	}
	resource := &Resource{
		Cluster:         cluster,
		OwnerUID:        ownerUID,
		UID:             accessor.GetUID(),
		Name:            accessor.GetName(),
		Namespace:       accessor.GetNamespace(),
		Group:           f.storageGroupResource.Group,
		Resource:        f.storageGroupResource.Resource,
		Version:         f.storageVersion.Version,
		Kind:            obj.GetObjectKind().GroupVersionKind().Kind,
		ResourceVersion: accessor.GetResourceVersion(),
		Object:          buf.Bytes(),
		CreatedAt:       accessor.GetCreationTimestamp().Time,
	}
	f.DB.Table = append(f.DB.Table, resource)
	index := len(f.DB.Table) - 1
	f.DB.Index[Cluster][cluster] = append(f.DB.Index[Cluster][cluster], index)
	f.DB.Index[Namespace][resource.Namespace] = append(f.DB.Index[Namespace][resource.Namespace], index)
	return nil
}

func (f *FakeResourceStorage) Update(ctx context.Context, cluster string, obj runtime.Object) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	index, err := f.GetIndex(ctx, cluster, accessor.GetNamespace(), accessor.GetName())
	if err != nil {
		return err
	}
	var ownerUID types.UID
	if owner := metav1.GetControllerOfNoCopy(accessor); owner != nil {
		ownerUID = owner.UID
	}
	var buf bytes.Buffer
	if err := f.codec.Encode(obj, &buf); err != nil {
		return err
	}
	f.DB.Table[index] = &Resource{
		Cluster:         cluster,
		OwnerUID:        ownerUID,
		UID:             accessor.GetUID(),
		Name:            accessor.GetName(),
		Namespace:       accessor.GetNamespace(),
		Group:           f.storageGroupResource.Group,
		Resource:        f.storageGroupResource.Resource,
		Version:         f.storageVersion.Version,
		Kind:            obj.GetObjectKind().GroupVersionKind().Kind,
		ResourceVersion: accessor.GetResourceVersion(),
		Object:          buf.Bytes(),
		CreatedAt:       accessor.GetCreationTimestamp().Time,
	}
	return nil
}

func (f *FakeResourceStorage) Delete(ctx context.Context, cluster string, obj runtime.Object) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	index, err := f.GetIndex(ctx, cluster, accessor.GetNamespace(), accessor.GetName())
	if err != nil {
		return err
	}
	f.DB.Table[index] = nil
	var targetc, targetn int
	for k, v := range f.DB.Index[Cluster][cluster] {
		if v == index {
			targetc = k
		}
	}
	for k, v := range f.DB.Index[Namespace][accessor.GetNamespace()] {
		if v == index {
			targetn = k
		}
	}
	var tmp []int
	tmp = f.DB.Index[Cluster][cluster][:targetc]
	tmp = append(tmp, f.DB.Index[Cluster][cluster][targetc+1:]...)
	f.DB.Index[Cluster][cluster] = tmp
	tmp = f.DB.Index[Namespace][accessor.GetNamespace()][:targetn]
	tmp = append(tmp, f.DB.Index[Namespace][accessor.GetNamespace()][targetn+1:]...)
	f.DB.Index[Namespace][accessor.GetNamespace()] = tmp
	return nil
}

func (f *FakeResourceStorage) GetIndex(ctx context.Context, cluster, namespace, name string) (int, error) {
	index := f.DB.UnionClusterNamespace(cluster, namespace)
	for _, v := range index {
		if f.DB.Table[v].Name == name {
			return v, nil
		}
	}
	return -1, utils.NewError("can not find index for this resource")
}

// DataBase
// Table similar to table in mysql for storing data
// Index can be thought of as an index similar to mysql
// and Slices in map must be sequential
type DataBase struct {
	Table    []*Resource
	Index    []map[string][]int
	FreeList []int
}

type IndexType int

const (
	Id IndexType = iota
	Resources
	Cluster
	Namespace
	ClusterNamespace
)

func NewDataBase() *DataBase {
	m := make([]map[string][]int, ClusterNamespace)
	m[Resources] = make(map[string][]int)
	m[Cluster] = make(map[string][]int)
	m[Namespace] = make(map[string][]int)
	return &DataBase{
		Table: make([]*Resource, 0),
		Index: m,
	}
}

func (db *DataBase) UnionClusterNamespace(cluster, namespace string) []int {
	var res []int
	var ci, ni, clen, nlen int
	clen = len(db.Index[Cluster][cluster])
	nlen = len(db.Index[Namespace][namespace])
	for ci < clen && ni < nlen {
		if db.Index[Cluster][cluster][ci] < db.Index[Namespace][namespace][ni] {
			ci++
		} else if db.Index[Cluster][cluster][ci] > db.Index[Namespace][namespace][ni] {
			ni++
		} else {
			res = append(res, db.Index[Cluster][cluster][ci])
			ci++
			ni++
		}
	}
	return res
}

func (db *DataBase) ListSimpleSearch(opts *internal.ListOptions) []*Resource {
	if opts == nil {
		return db.Table
	}
	var res []*Resource
	var index []int
	if opts.ClusterNames == nil {
		opts.ClusterNames = db.GetAllCluster()
	}
	if opts.Namespaces == nil {
		opts.Namespaces = db.GetAllNamespace()
	}
	for _, cluster := range opts.ClusterNames {
		for _, namespace := range opts.Namespaces {
			index = append(index, db.UnionClusterNamespace(cluster, namespace)...)
		}
	}

	for _, v := range index {
		res = append(res, db.Table[v])
	}
	if opts.ExtraLabelSelector != nil {
		owner := db.OwnerId(opts.ExtraLabelSelector)
		if owner != nil {
			res = db.FindOwner(res, owner)
		}

		tmp, err := db.LabelSelect(res, opts.ExtraLabelSelector)
		if err != nil {
			return nil
		}
		res = tmp
	}
	if opts.LabelSelector != nil {
		owner := db.OwnerId(opts.LabelSelector)
		if owner != nil {
			res = db.FindOwner(res, owner)
		}

		tmp, err := db.LabelSelect(res, opts.LabelSelector)
		if err != nil {
			return nil
		}
		res = tmp
	}
	if opts.EnhancedFieldSelector != nil && opts.EnhancedFieldSelector.String() != "" {
		tmp, err := db.EnhancedSelect(res, opts.EnhancedFieldSelector)
		if err != nil {
			return nil
		}
		res = tmp
	}
	if opts.OwnerName != "" {
		count := opts.OwnerSeniority
		var ownerlist []*OwnerSearch
		for _,v := range res {
			s := &OwnerSearch{
				Self: v,
				Owner: []*Resource{v},
			}
			ownerlist = append(ownerlist, s)
		}
		for count >= 0 {
			if count != 0 {
				var err error
				ownerlist, err = db.FindOwnerByName(ownerlist)
				if err != nil {
					break
				}
			} else {
				tmp, err := db.OwnerName(ownerlist, opts.OwnerName)
				if err != nil {
					return nil
				}
				res = tmp
			}
			count --
		}
	}
	if opts.Limit > 0 {
		if opts.Limit > int64(len(res)) {
			return res
		}
		offset, err := strconv.Atoi(opts.Continue)
		if err != nil {
			offset = 0
		}
		if offset < 0 {
			return nil
		}
		start := offset
		end := int64(offset) + opts.Limit
		if start > len(res) {
			return nil
		} else if end > int64(len(res)) {
			return res[start:]
		} else {
			return res[start:end]
		}
	} else {
		return res
	}
}

func (db *DataBase) LabelSelect(list []*Resource, s labels.Selector) ([]*Resource, error) {
	var res []*Resource
	r, _ := s.Requirements()
	for _, v := range list {
		s2, err := v.Object.MarshalJSON()
		if err != nil {
			return nil, err
		}
		var m map[string]interface{}
		err = json.Unmarshal(s2, &m)
		if err != nil {
			return nil, err
		}
		maps, ok := m["metadata"].(map[string]interface{})["labels"]
		if ok {
			sl := maps.(map[string]interface{})
			for key, value := range sl {
				for _, requirement := range r {
					if key == requirement.Key() {
						for index, _ := range requirement.Values() {
							if value == index {
								res = append(res, v)
								continue
							}
						}
					}
				}
			}
		}
	}
	return res, nil
}

func (db *DataBase) EnhancedSelect(list []*Resource, s fields2.Selector) ([]*Resource, error) {
	var res []*Resource
	if r, b := s.Requirements(); b {
		for _, obj := range list {
			s2, err := obj.Object.MarshalJSON()
			if err != nil {
				return nil, err
			}
			var target map[string]interface{}
			err = json.Unmarshal(s2, &target)
			if err != nil {
				return nil, err
			}
			for _, v := range r {
				targetList := strings.Split(v.Fields()[len(v.Fields())-1].Path().String(), ".")
				for index, field := range targetList {
					if index < len(targetList)-1 {
						var ok bool
						target, ok = target[field].(map[string]interface{})
						if !ok {
							return nil, err
						}
					} else {
						n, ok := target[field].(string)
						if !ok {
							return nil, err
						}
						for _, value := range v.Values().List() {
							if value == n {
								res = append(res, obj)
							}
						}
					}
				}
			}
		}
	}
	return res, nil
}

func (db *DataBase) OwnerId(s labels.Selector) []string {
	var res []string
	r, _ := s.Requirements()
	for _, v := range r {
		if v.Key() == constants.SearchLabelOwnerUID {
			for key, _ := range v.Values() {
				res = append(res, key)
			}
			return res
		}
	}
	return nil
}

func (db *DataBase) OwnerName(list []*OwnerSearch, name string) ([]*Resource, error) {
	var res []*Resource
	for _,os := range list {
		for _, obj := range os.Owner {
			bt, err := obj.Object.MarshalJSON()
			if err != nil {
				return nil, err
			}
			var target map[string]interface{}
			err = json.Unmarshal(bt, &target)
			if err != nil {
				return nil, err
			}
			if metadata, ok := target["metadata"].(map[string]interface{}); ok {
				if owners, ok := metadata["ownerReferences"].([]interface{}); ok {
					for _, ref := range owners {
						body, ok := ref.(map[string]interface{})
						if ok && body["name"] == name {
							res = append(res, obj)
						}
					}
				}
			}
		}
	}
	return res, nil
}

func (db *DataBase) FindOwnerByName(list []*OwnerSearch) ([]*OwnerSearch, error) {
	for _,v := range list {
		var owner []*Resource
		for _,value := range v.Owner {
			bt,err := value.Object.MarshalJSON()
			if err != nil {
				return nil,err
			}
			var target map[string]interface{}
			err = json.Unmarshal(bt, &target)
			if err != nil {
				return nil, err
			}
			if metadata, ok := target["metadata"].(map[string]interface{}); ok {
				if owners, ok := metadata["ownerReferences"].([]interface{}); ok {
					for _, ref := range owners {
						body, ok := ref.(map[string]interface{})
						if ok {
							owner = append(owner, db.FindResourceByName(body["name"].(string)))
						}
					}
				}
			}

		}
		v.Owner = owner
	}
	return list,nil
}

func (db *DataBase) FindResourceByName(name string) *Resource {
	for _,v := range db.Table {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (db *DataBase) FindOwner(res []*Resource, owner []string) []*Resource {
	var tmp []*Resource
	for _, v := range res {
		for _, id := range owner {
			if string(v.OwnerUID) == id {
				tmp = append(tmp, v)
			}
		}
	}
	if tmp == nil {
		return res
	}
	return tmp
}

func (db *DataBase) GetAllCluster() []string {
	var res []string
	for key, _ := range db.Index[Cluster] {
		res = append(res, key)
	}
	return res
}

func (db *DataBase) GetAllNamespace() []string {
	var res []string
	for key, _ := range db.Index[Namespace] {
		res = append(res, key)
	}
	return res
}

type OwnerSearch struct {
	Self	*Resource
	Owner	[]*Resource
}