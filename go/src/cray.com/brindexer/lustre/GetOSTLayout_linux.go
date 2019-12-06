// +build linux

package lustre

/*
#include <stdint.h>
#include <stdlib.h>
#include <sys/types.h>
#include <linux/lustre/lustre_user.h>
#include <linux/lustre/lustre_idl.h>
typedef struct {
	int			um_numosts;
	enum lov_comp_md_flags	um_mirror_state;
	uint32_t		um_mdtidx;
	uint32_t		um_ostidx[LOV_MAX_STRIPE_COUNT];
    char* 			um_pool_name;
}C_user_layout, *pC_user_layout;

static inline struct lov_user_md *
lov_comp_v1_entry(struct lov_comp_md_v1 *comp_v1, int ent_idx)
{
	return (struct lov_user_md *)((char *)comp_v1 +
			comp_v1->lcm_entries[ent_idx].lcme_offset);
}

static inline struct lov_user_ost_data_v1 *
lov_v1v3_objects(struct lov_user_md *v1)
{
	if (v1->lmm_magic == LOV_USER_MAGIC_V3)
		return ((struct lov_user_md_v3 *)v1)->lmm_objects;
	else
		return v1->lmm_objects;
}
void free_info(pC_user_layout info){
        free(info->um_pool_name);
        free(info);
}

pC_user_layout lustre_get_layout_file(const char *path)
{
	ssize_t size=0;
	char buf[XATTR_SIZE_MAX] = {0};
	struct lov_user_md_v1 *lum=0;
	int rc = 0;
	int i = 0;
	pC_user_layout info=0;

	size = getxattr(path, XATTR_NAME_LOV, buf, sizeof(buf));
	if (size < 0)
		return 0;
	info = (pC_user_layout)malloc(sizeof(C_user_layout));
	memset(info, 0, sizeof(C_user_layout));
	info->um_pool_name = (char*)malloc(LOV_MAXPOOLNAME+1);
	memset(info->um_pool_name, 0, LOV_MAXPOOLNAME+1);
	lum = (struct lov_user_md *)buf;

	switch (lum->lmm_magic) {
	case LOV_USER_MAGIC_V3: {
		struct lov_user_md_v3 *v3 = (struct lov_user_md_v3 *)buf;

		snprintf(info->um_pool_name, LOV_MAXPOOLNAME+1, "%s",
			 v3->lmm_pool_name);
	}
	case LOV_USER_MAGIC_V1:
	case LOV_USER_MAGIC_SPECIFIC:
		for (i = 0; i < lum->lmm_stripe_count; i++)
			info->um_ostidx[i] = lum->lmm_objects[i].l_ost_idx;

		info->um_numosts = (int)lum->lmm_stripe_count;

		break;

	case LOV_USER_MAGIC_COMP_V1: {
		struct lov_comp_md_v1 *comp_v1 = (struct lov_comp_md_v1 *)lum;
		struct lov_comp_md_entry_v1 *entry;
		struct lov_user_md_v1 *v1;
		struct lov_user_ost_data_v1 *objects;
		int j=0;
		int k=0;

		for (i = 0, k = 0; i < comp_v1->lcm_entry_count; i++) {
			entry = &comp_v1->lcm_entries[i];
			if (!(entry->lcme_flags & LCME_FL_INIT))
				continue;

			v1 = lov_comp_v1_entry(comp_v1, i);
			objects = lov_v1v3_objects(v1);

			for (j = 0; j < v1->lmm_stripe_count; j++, k++)
				info->um_ostidx[k] = objects[j].l_ost_idx;

			info->um_numosts += (int)v1->lmm_stripe_count;
		}

		info->um_mirror_state = comp_v1->lcm_flags & LCM_FL_FLR_MASK;
	}
		break;

	case LMV_MAGIC_V1: {
		struct lmv_mds_md_v1 *lmm = (struct lmv_mds_md_v1 *)buf;

		snprintf(info->um_pool_name, LOV_MAXPOOLNAME+1 , "%s",
			 lmm->lmv_pool_name);

		info->um_mdtidx = lmm->lmv_master_mdt_index;
	}
		break;

	default:
		//fprintf(stderr, "Unknown file magic %u\n", lum->lmm_magic);
		free_info(info);
		return 0;
	}
	return info;
}



*/
import "C"
import (
	"fmt"
	"unsafe"
)

func GetLayout(path string) *OSTLayout {

	cpath := C.CString(path) //*C.char
	defer C.free(unsafe.Pointer(cpath))
	clayout := C.lustre_get_layout_file(cpath)

	osts := make([]uint32, 0, 64)
	layout := OSTLayout{0, "", osts}
	if clayout == nil {
		fmt.Println("lustre_get_layout_file failed:", path)
		return &layout
	}
	defer C.free(unsafe.Pointer(clayout))
	defer C.free(unsafe.Pointer(clayout.um_pool_name))

	//Convert
	layout.mirrorState = clayout.um_mirror_state
	layout.poolName = C.GoString(clayout.um_pool_name)
	for i := C.int(0); i < clayout.um_numosts; i++ {
		layout.ostIndice = append(layout.ostIndice, uint32(clayout.um_ostidx[i]))
	}
	return &layout

}
