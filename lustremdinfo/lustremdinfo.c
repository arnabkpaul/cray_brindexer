/*
 * Lustre indexing library for Brindexer
 * Copyright 2019 Cray Inc. All Rights Reserved.
 */

#include <stdio.h>
#include <string.h>
#include <limits.h>
#include <errno.h>
#include <fcntl.h>
#include <sys/xattr.h>
#include <sys/types.h>
#include <sys/stat.h>

#include <lustre/lustre_user.h>
#include <linux/lustre/lustre_idl.h>

#include "lustremdinfo.h"




/**
 * Obtain md information for a file
 *
 * Users of the API should define a struct user_mdinfo instance, set its memory
 * to 0 and pass its address as the info parameter. If statx information is
 * required, the all parameter should be set to true
 *
 * \param path The file's absolute path
 * \param all  When true, obtain all file information in struct user_mdinfo;
 *	       when false, only obtain statx information
 * \param info The md information for the file. This function allocates it and
 *	       the caller needs to free it
 *
 * \retval 0 for success
 * \retval -ve on failure
 */

int lustre_get_info_file(const char *path, bool all, struct user_mdinfo *info)
{
	ssize_t size;
	char buf[XATTR_SIZE_MAX] = {0};
	struct lov_user_md_v1 *lum;
	int i;
	int rc = 0;

	if (!path || !info)
		return -EINVAL;

	size = getxattr(path, XATTR_NAME_LOV, buf, sizeof(buf));
	if (size < 0)
		return -errno;

	lum = (struct lov_user_md *)buf;

	switch (lum->lmm_magic) {
	case LOV_USER_MAGIC_V3: {
		struct lov_user_md_v3 *v3 = (struct lov_user_md_v3 *)buf;

		/* Pool name */
		snprintf(info->um_pool_name, sizeof(info->um_pool_name), "%s",
			 v3->lmm_pool_name);

		/* fallthrough */
	}
	case LOV_USER_MAGIC_V1:
	case LOV_USER_MAGIC_SPECIFIC:
		/* OST list for plain layouts */
		for (i = 0; i < lum->lmm_stripe_count; i++)
			info->um_ostidx[i] = lum->lmm_objects[i].l_ost_idx;

		info->um_numosts = lum->lmm_stripe_count;

		break;

	case LOV_USER_MAGIC_COMP_V1: {
		struct lov_comp_md_v1 *comp_v1 = (struct lov_comp_md_v1 *)lum;
		struct lov_comp_md_entry_v1 *entry;
		struct lov_user_md_v1 *v1;
		struct lov_user_ost_data_v1 *objects;
		int i;
		int j;
		int k;

		/* OST list for composite layouts */
		for (i = 0, k = 0; i < comp_v1->lcm_entry_count; i++) {
			entry = &comp_v1->lcm_entries[i];

			/* Non-instantianted sub-layouts don't have any
			 * objects
			 */
			if (!(entry->lcme_flags & LCME_FL_INIT))
				continue;

			v1 = (struct lov_user_md *)((char *)comp_v1 +
			      comp_v1->lcm_entries[i].lcme_offset);
			objects = (v1->lmm_magic == LOV_USER_MAGIC_V3) ?
				  ((struct lov_user_md_v3 *)v1)->lmm_objects :
				  v1->lmm_objects;

			for (j = 0; j < v1->lmm_stripe_count; j++, k++)
				info->um_ostidx[k] = objects[j].l_ost_idx;

			info->um_numosts += v1->lmm_stripe_count;
		}

		/* Mirror state */
		info->um_mirror_state = comp_v1->lcm_flags & LCM_FL_FLR_MASK;
	}
		break;

	case LMV_MAGIC_V1: {
		struct lmv_mds_md_v1 *lmm = (struct lmv_mds_md_v1 *)buf;

		/* Pool name */
		snprintf(info->um_pool_name, sizeof(info->um_pool_name), "%s",
			 lmm->lmv_pool_name);

		/* MDT index */
		info->um_mdtidx = lmm->lmv_master_mdt_index;
	}
		break;

	default:
		fprintf(stderr, "Unknown file magic %u\n", lum->lmm_magic);

		return -EINVAL;
	}

	/* HSM attributes */
	memset(buf, 0, sizeof(buf));

	/* TODO: These xattrs should probably be swabbed for endianness; e.g.
	 * lustre_hsm_swab(); this is not an issue on x86_64
	 */
	size = getxattr(path, XATTR_NAME_HSM, buf, sizeof(buf));
	if (size < 0) {
		/* Ignore non-existent HSM xattrs */
		if (errno != ENODATA)
			return -errno;
	} else {
		memcpy(&info->um_hsm_attrs, buf, size);
	}

	if (!all)
		goto out;

	/* stat */
	rc = stat(path, &info->um_stat);
	if (rc)
		rc = -errno;

out:

	return rc;
}
