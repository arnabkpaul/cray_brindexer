/*
 * Lustre indexing library for Brindexer
 * Copyright 2019 Cray Inc. All Rights Reserved.
 */
/*
 * This program demonstrates an example usage of the lustre_get_info_file() API
 * call that prints out some statx information and the number and indexes of
 * OSTs for both plain and composite layout types.
 *
 * It can be built by issuing "make test"
 */
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <lustre/lustreapi.h>
#include <lustre/lustre_user.h>

#include "lustremdinfo.h"


int main(int argc, char **argv)
{
	struct user_mdinfo info = {0};
	int root_fd;
	int i;
	int rc;


	if (argc < 2) {
		fprintf(stderr, "Error: %s, exiting!\n", strerror(EINVAL));
		return -EINVAL;
	}

	memset(&info, 0, sizeof(info));

	rc = lustre_get_info_file(argv[1], true, &info);
	if (rc) {
		fprintf(stderr, "Error: %s, exiting!\n", strerror(rc));
		goto out;
	}

	printf("st_ino = %llu, st_size = %llu, st_uid = %lu, "
	       "st_gid = %lu\nnum_osts = %u\n", info.um_stat.st_ino,
	       info.um_stat.st_size, info.um_stat.st_uid,
	       info.um_stat.st_gid, info.um_numosts);

	for (i = 0; i < info.um_numosts; i++)
		printf("ost idx = %u\n", info.um_ostidx[i]);

out:
	return rc;
}
